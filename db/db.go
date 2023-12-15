package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Item struct {
	ID             int      `yaml:"id"`
	Name           string   `yaml:"name"`
	EQTC           int      `yaml:"eqtc"`
	Alla           int      `yaml:"alla"`
	SoldBy         []string `yaml:"sold_by"`
	QuestRewarded  []string `yaml:"quest_rewarded"`
	QuestReagent   []string `yaml:"quest_reagent"`
	DroppedBy      []string `yaml:"dropped_by"`
	RecipeRewarded []string `yaml:"recipe_rewarded"`
	RecipeReagent  []string `yaml:"recipe_reagent"`
}

type Merchant struct {
	Name  string `yaml:"name"`
	Loc   string `yaml:"loc"`
	Items []int  `yaml:"items"`
}

type Npc struct {
	Name string `yaml:"name"`
	Min  int    `yaml:"min"`
	Max  int    `yaml:"max"`
	Alla int    `yaml:"alla"`
}

type Quest struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	Alla int    `yaml:"alla"`
	Min  int    `yaml:"min"`
	Max  int    `yaml:"max"`
}

type Recipe struct {
	ID         int               `yaml:"id"`
	Name       string            `yaml:"name"`
	Trivial    int               `yaml:"trivial"`
	RecipeType string            `yaml:"recipe_type"`
	OtherID    int               `yaml:"other_id"`
	ResultID   int               `yaml:"result_id"`
	Components []RecipeComponent `yaml:"components"`
}

type RecipeComponent struct {
	ComponentID int `yaml:"component_id"`
	ItemID      int `yaml:"item_id"`
	OtherID     int `yaml:"other_id"`
	OtherType   int `yaml:"other_type"`
}

func ParseItems(items map[int]*Item) error {
	err := filepath.WalkDir("db/item/", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}
		if filepath.Base(path) == "_base.yaml" {
			return nil
		}
		r, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer r.Close()
		type itemStruct struct {
			Entry []Item `yaml:"items"`
		}
		itemFile := &itemStruct{}
		err = yaml.NewDecoder(r).Decode(itemFile)
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		for i := range itemFile.Entry {
			item := itemFile.Entry[i]
			for j := range item.DroppedBy {
				item.DroppedBy[j] = strings.ReplaceAll(item.DroppedBy[j], " ", "_")
			}
			for j := range item.SoldBy {
				item.SoldBy[j] = strings.ReplaceAll(item.SoldBy[j], " ", "_")
			}
			for j := range item.QuestReagent {
				item.QuestReagent[j] = strings.ReplaceAll(item.QuestReagent[j], " ", "_")
			}
			for j := range item.QuestRewarded {
				item.QuestRewarded[j] = strings.ReplaceAll(item.QuestRewarded[j], " ", "_")
			}

			_, ok := items[item.ID]
			if ok {
				return fmt.Errorf("duplicate item id: %d", item.ID)
			}

			items[item.ID] = &item
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}
	return nil
}
