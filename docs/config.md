# Configuration

The configuration is located in a file named `config.ini` in the current working directory.

The file is created/updated automatically with default values when missing.

## general

### online

**Whether or not to use wallhaven's API.**

```ini
online = true
```

When enabled, it requires an internet connection to search and download wallpapers.

When false, it does not require an internet connection, and the previously kept wallpapers are reused instead.

Small notes:

* The online mode fallback on the offline mode if there are connectivity issues.
* The offline mode does nothing if no wallpapers were kept.

### keep

**Whether or not to keep the wallpapers that were downloaded.**

```ini
keep = false
```

This is useful for use with the offline mode or perhaps you like hoarding data.

### directory

**Location for the wallpapers**

```ini
directory = "downloads"
```

The offline mode uses wallpapers from that directory.
 
Downloaded wallpapers also saved in that directory.

### interval

**Rotate the wallpaper at a given interval**

```ini
interval = 0
```

The interval is written in a human-readable format.

The notation uses `y = year`, `mo = month`, `w = week`, `d = days`, `h = hours`, `m = minutes` and `s = seconds`.

For example, `2w5d` is 2 weeks and 5 days.

When a non-zero interval is configured, the wallpaper updates immediately, then the program runs in the background to
continue updating it every interval.

When an interval of 0 is configured, the wallpaper updates once then the program exits immediately.

## randomize

```ini
randomize = true
```

Whether or not to randomize the wallpaper selection from the depth pool available.

Otherwise, it'll go through them sequentially.

# search

## tags

Comma-separated tags to search.

```ini
tags = "fantasy"
```

## categories

```ini
general = true
anime = true
people = true
```

## sorting

```ini
sorting = "favorites"
```

Must be one of `"dateAdded"`, `"relevance"`, `"random"`, `"views"`, `"favorites"` or `"topList"`.

## purities

```ini
sfw = true
sketchy = true
nsfw = false
```

NSFW requires an API token and not supported yet.

## resolution

**Resolution for the wallpapers**

```ini
resolution = 1920x1080
```

## exact

```ini
exact = false
```
Treat the resolution as a minimum requirement, or force the wallpapers to be exactly what's configured.

## ratio

**Aspect-ratio for the wallpapers**

```ini
ratio = 16:9
```

## depth

**How many wallpapers to consider from a given search result**

Some searches results into thousands of wallpapers. This setting configures
how many wallpapers to take into account.  The remaining ones are ignored.

This is necessary because you may want to cycle through only the top-100
most favorited wallpapers, without going too deep, as the deeper you go, the
worse the quality of the results become.

A depth of 0 uses all the results.

```ini
depth = 100
```

## colors

```ini
colors = ""
```

Wallpapers with a strong preference for the specified colors, comma-separated.

Empty string if you have no preference.
