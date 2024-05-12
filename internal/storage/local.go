package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"

	"github.com/akl-infra/slf/v2"
	"github.com/charmbracelet/log"
	"golang.org/x/exp/maps"
)

var Path string
var Cache map[string]slf.Layout = make(map[string]slf.Layout)

func readLayout(name string) (slf.Layout, error) {
	var slfLayout slf.Layout
	data, err := os.ReadFile(Path + name)
	if err != nil {
		return slf.Layout{}, err
	}
	err = json.Unmarshal(data, &slfLayout)
	if err != nil {
		return slf.Layout{}, err
	}
	return slfLayout, nil
}

func writeLayout(layout slf.Layout) error {
	data, err := json.Marshal(layout)
	if err != nil {
		return err
	}
	err = os.WriteFile(Path+layout.Name, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Init(path string) error {
	if strings.HasSuffix(path, "/") {
		Path = path
	} else {
		Path = path + "/"
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Error(err)
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			name := entry.Name()
			layout, err := readLayout(name)
			if err != nil {
				log.Error(err)
				return err
			}
			Cache[name] = layout
		}
	}

	return nil
}

func Get(name string) (slf.Layout, error) {
	if layout, ok := Cache[name]; ok {
		return layout, nil
	} else {
		return slf.Layout{}, errors.New("Layout not found")
	}
}

func Put(layout slf.Layout) error {
	Cache[layout.Name] = layout
	err := writeLayout(layout)

	return err
}

func List() []string {
	layouts := maps.Keys(Cache)
	sort.Sort(sort.StringSlice(layouts))
	return layouts
}
