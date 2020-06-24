// -ldflags -H=windowsgui

package main

import (
	"github.com/dlasky/go-wallhaven"
	"github.com/reujab/wallpaper"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var cachedWallpapers []wallhaven.Wallpaper
var cachedWallpaperIds map[wallhaven.WallpaperID]bool
var cachedWallpapersMutex sync.Mutex
var cachedWallpapersOnce sync.Once

var nextIndex int

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	handleDragAndDrop()
	handleRegularUse()
}

// This is a special case that gives the user the ability to drag-and-drop files onto the executable to
// set a new wallpaper.
func handleDragAndDrop() {
	if len(os.Args) > 1 && os.Args[1] != "" {
		_ = wallpaper.SetFromFile(os.Args[1])
		os.Exit(0)
	}
}

func handleRegularUse() {
	for {
		config, err := loadConfig()
		if err != nil {
			log.Fatalln(err)
		}

		// The online mode contacts the api, while the offline mode re-uses previously downloaded wallpapers.
		if config.Online {
			err = updateWallpaperOnline(config)
		}

		// If there was an error with the online mode, fallback on offline mode.
		// Sometimes people want to explicitly use the offline mode, which is fine too.
		if err != nil || !config.Online {
			err = updateWallpaperOffline(config)
		}

		// At this point, any error is unrecoverable.
		if err != nil {
			panic(err)
		}

		// No need to run in a loop when there's no interval. Simply exit.
		if config.Interval == 0 {
			os.Exit(0)
		}

		time.Sleep(config.Interval)
	}
}

func updateWallpaperOnline(config *Config) error {
	search := wallhaven.Search{
		Categories: config.Categories,
		Purities: config.Purities,
		Sorting: config.Sorting,
		Order: wallhaven.Desc,
		AtLeast: config.Resolution,
		Colors: config.Colors,
		Ratios: []wallhaven.Ratio{
			{
				Horizontal: config.Ratio.Horizontal,
				Vertical: config.Ratio.Vertical,
			},
		},
		Query: wallhaven.Q{
			Tags: config.Tags,
		},
	}

	// Must wait for at least 1 wallpaper to proceed, but we can let the rest get cached in the background.
	cachedWallpapersOnce.Do(func() {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go cachingWallpapers(config, search, &wg)
		wg.Wait()
	})

	// TODO: Support a linear index too

	cachedWallpapersMutex.Lock()
	defer cachedWallpapersMutex.Unlock()

	if config.Randomize {
		nextIndex = rand.Intn(len(cachedWallpapers) - 1)
	}

	// Save to disk if enabled.
	if config.Keep {
		err := downloadImageURL(config, cachedWallpapers[nextIndex].Path)
		if err != nil {
			return err
		}

		return wallpaper.SetFromFile(cachedWallpapers[nextIndex].Path)
	}

	err := wallpaper.SetFromURL(cachedWallpapers[nextIndex].Path)
	if err != nil {
		return err
	}

	nextIndex++
	if nextIndex == len(cachedWallpapers) {
		nextIndex = 0
	}

	return nil
}

func cachingWallpapers(config *Config, search wallhaven.Search, wg *sync.WaitGroup) {
	search.Page = 1

	for {
		results, err := wallhaven.SearchWallpapers(&search)
		if err != nil {
			// Retry on error, after waiting a little bit.
			time.Sleep(time.Minute)
			continue
		}

		for _, result := range results.Data {
			cachedWallpapersMutex.Lock()

			if cachedWallpaperIds == nil {
				cachedWallpaperIds = make(map[wallhaven.WallpaperID]bool)
			}

			// Check if we already have the wallpaper.
			found := cachedWallpaperIds[result.ID]

			// If not, add it.
			if !found {
				if len(cachedWallpapers) == 0 {
					wg.Done()
				}

				cachedWallpaperIds[result.ID] = true
				cachedWallpapers = append(cachedWallpapers, result)
			}

			cachedWallpapersMutex.Unlock()
		}

		cachedWallpapersMutex.Lock()
		if len(cachedWallpapers) > config.Depth {
			cachedWallpapersMutex.Unlock()
			break
		}
		cachedWallpapersMutex.Unlock()

		if results.Meta.CurrentPage == results.Meta.LastPage {
			break
		}

		search.Page++

		// Must wait in-between requests, the API allows 45 requests per minute and is very sensitive to
		// bursts of requests.
		time.Sleep(time.Duration(int64(math.Max(float64(5 * time.Second), float64(config.Interval)))))
	}
}

func downloadImageURL(config *Config, url string) error {
	err := os.MkdirAll(config.Directory, 0700)
	if err != nil {
		return err
	}

	response, err := http.Get(url)
	if err != nil {
		return err
	}

	filePath := filepath.Join(config.Directory, filepath.Base(url))
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func updateWallpaperOffline(config *Config) error {
	_, err := os.Stat(config.Directory)

	if os.IsNotExist(err) {
		return err
	}

	var files []string

	err = filepath.Walk(config.Directory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	})

	if err != nil {
		return err
	}

	// It's fine if the directory does not exist, we simply do nothing.
	if len(files) == 0 {
		return nil
	}

	rand.Shuffle(len(files), func(i, j int) {
		files[i], files[j] = files[j], files[i]
	})

	absoluteFilepath, err := filepath.Abs(files[0])
	if err != nil {
		return err
	}

	return wallpaper.SetFromFile(absoluteFilepath)
}