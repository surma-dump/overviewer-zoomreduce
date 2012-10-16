package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/voxelbrain/goptions"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	VERSION     = "0.1.1"
	CONFIG_FILE = "overviewerConfig.js"
)

var (
	options = struct {
		goptions.Help `goptions:"-h, --help"`
		goptions.Verbs
		List struct {
		} `goptions:"list"`
		Remove struct {
			Num    int      `goptions:"-n, --num, description='Number of zoom levels to remove'"`
			Worlds []string `goptions:"-w, --world, description='Include world folder in operation'"`
		} `goptions:"remove"`
	}{}
)

func init() {
	goptions.ParseAndFail(&options)
}

func main() {
	config, err := parseConfig()
	if err != nil {
		log.Fatalf("Could not parse config: %s", err)
	}

	switch options.Verbs {
	case "list":
		fmt.Println("Worlds:")
		for _, worldpath := range config.Worlds() {
			world := config.World(worldpath)
			fmt.Printf("\t%s (Depth: %d)\n", worldpath, world.ZoomLevels())
		}
	case "remove":
		for _, worldname := range options.Remove.Worlds {
			world := config.World(worldname)
			newzoom := world.ZoomLevels() - options.Remove.Num
			filepath.Walk(world.Path(), func(path string, fi os.FileInfo, err error) error {
				if fi.IsDir() {
					return nil
				}
				pe := strings.Split(path, "/")
				zoomLevel := len(pe)
				if zoomLevel > newzoom {
					log.Printf("Removing %s...", path)
					err := os.Remove(path)
					if err != nil {
						log.Printf("Could not delete %s: %s", path, err)
					}
				}
				return nil
			})
			world.SetZoomLevels(newzoom)
		}
		f, err := os.Create(CONFIG_FILE)
		if err != nil {
			log.Fatalf("Could not open config file %s for writing: %s", CONFIG_FILE, err)
		}
		defer f.Close()
		io.WriteString(f, "var overviewerConfig = ")
		enc := json.NewEncoder(f)
		enc.Encode(config)
	default:
		goptions.PrintHelp()
		return
	}
}

type Config map[string]interface{}

func (c Config) Worlds() []string {
	worlds := c["tilesets"].([]interface{})
	r := make([]string, 0, len(worlds))
	for _, world := range worlds {
		r = append(r, world.(map[string]interface{})["path"].(string))
	}
	return r
}

type World map[string]interface{}

func (c Config) World(worldpath string) World {
	for _, tileset := range c["tilesets"].([]interface{}) {
		ts := tileset.(map[string]interface{})
		if ts["path"].(string) == worldpath {
			return ts
		}
	}
	panic("Invalid world name")
}

func (w World) Path() string {
	return w["path"].(string)
}

func (w World) ZoomLevels() int {
	return int(w["maxZoom"].(float64))

}

func (w World) SetZoomLevels(zl int) {
	w["maxZoom"] = float64(zl)
}

func parseConfig() (Config, error) {
	f, err := os.Open(CONFIG_FILE)
	if err != nil {
		log.Fatalf("Could not open overviewerConfig.js: %s", err)
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for r := rune(' '); r != rune('{'); {
		var err error
		r, _, err = br.ReadRune()
		if err != nil {
			return nil, err
		}
	}
	br.UnreadRune()

	c := make(Config)
	dec := json.NewDecoder(br)
	return c, dec.Decode(&c)
}
