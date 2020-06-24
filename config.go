package main

import (
	"bytes"
	"fmt"
	"github.com/dlasky/go-wallhaven"
	"github.com/go-ini/ini"
	"github.com/karrick/tparse"
	"os"
	"strings"
	"time"
)

const ConfigLocation = "config.ini"

type Config struct {
	// General
	Online bool
	Keep bool
	Directory string
	Interval time.Duration
	Randomize bool

	// Search
	Tags []string
	Sorting wallhaven.Sort
	Categories wallhaven.Category
	Purities wallhaven.Purity
	Resolution wallhaven.Resolution
	Ratio wallhaven.Ratio
	Exact bool
	Depth int
	Colors []string
}

func createConfigIfNecessary() error {
	_, err := os.Stat(ConfigLocation)

	if os.IsNotExist(err) {
		file, err := os.Create(ConfigLocation)
		if err != nil {
			return err
		}

		cleanDefaultConfig := strings.ReplaceAll(strings.TrimSpace(defaultConfig), "\t", "")

		buffer := bytes.Buffer{}
		buffer.WriteString(cleanDefaultConfig)
		_, err = buffer.WriteTo(file)
		if err != nil {
			return err
		}

		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

// TODO: Finish up the error checking. I got lazy.
func loadConfig() (*Config, error) {
	err := createConfigIfNecessary()
	if err != nil {
		return nil, err
	}

	cfg, err := ini.Load(ConfigLocation)
	if err != nil {
		return nil, err
	}

	config := Config{}

	config.Online, _ = cfg.Section("general").Key("online").Bool()
	config.Keep, _ = cfg.Section("general").Key("keep").Bool()
	config.Randomize, _ = cfg.Section("general").Key("randomize").Bool()

	config.Directory = cfg.Section("general").Key("directory").String()

	interval := cfg.Section("general").Key("interval").String()
	if interval != "0" {
		future, _ := tparse.ParseNow(time.RFC3339, "now+"+interval)
		config.Interval = time.Until(future)
	}

	tags := cfg.Section("search").Key("tags").String()
	config.Tags = strings.Split(tags, ",")

	sorting := cfg.Section("search").Key("sorting").String()
	switch sorting {
		case "dateAdded": config.Sorting = wallhaven.DateAdded
		case "relevance": config.Sorting = wallhaven.Relevance
		case "random": config.Sorting = wallhaven.Random
		case "views": config.Sorting = wallhaven.Views
		case "favorites": config.Sorting = wallhaven.Favorites
		case "topList": config.Sorting = wallhaven.Toplist
	}

	general, _ := cfg.Section("search").Key("general").Bool()
	anime, _ := cfg.Section("search").Key("anime").Bool()
	people, _ := cfg.Section("search").Key("people").Bool()

	config.Categories = 0

	if general {
		config.Categories &= wallhaven.General
	}

	if anime {
		config.Categories &= wallhaven.Anime
	}

	if people {
		config.Categories &= wallhaven.Anime
	}

	config.Purities = 0

	sfw, _ := cfg.Section("search").Key("sfw").Bool()
	sketchy, _ := cfg.Section("search").Key("sketchy").Bool()
	nsfw, _ := cfg.Section("search").Key("nsfw").Bool()

	if sfw {
		config.Purities &= wallhaven.SFW
	}

	if sketchy {
		config.Purities &= wallhaven.Sketchy
	}

	if nsfw {
		config.Purities &= wallhaven.NSFW
	}

	resolution := cfg.Section("search").Key("resolution").String()
	_, _ = fmt.Sscanf(resolution, "%dx%d", &config.Resolution.Width, &config.Resolution.Height)

	ratio := cfg.Section("search").Key("ratio").String()
	_, _ = fmt.Sscanf(ratio, "%d:%d", &config.Ratio.Horizontal, &config.Ratio.Vertical)

	config.Exact, _ = cfg.Section("search").Key("exact").Bool()
	config.Depth, _ = cfg.Section("search").Key("depth").Int()

	colors := cfg.Section("search").Key("colors").String()
	sliceColors := strings.Split(colors, ",")
	config.Colors = sliceColors

	return &config, nil
}

var defaultConfig = `
	[general]
	online = true
	keep = false
	directory = "downloads"
	interval = 0
	randomize = false
	
	[search]
	tags = ""
	
	general = true
	anime = true
	people = true
	
	colors = ""
	
	sorting = "favorites"
	depth = 100
	
	sfw = true
	sketchy = false
	nsfw = false
	
	resolution = 1920x1080
	ratio = 16:9
	exact = false
`