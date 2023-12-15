package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/codex/db"
	"gopkg.in/yaml.v3"
)

var (
	items     = make(map[int]*db.Item)
	merchants = make(map[string]*db.Merchant)
	npcs      = make(map[string]*db.Npc)
	quests    = make(map[string]*db.Quest)
	recipes   = make(map[int]*db.Recipe)
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed run: ", err)
		os.Exit(1)
	}
}

func run() error {
	err := db.ParseItems(items)
	if err != nil {
		return fmt.Errorf("parseItems: %w", err)
	}
	err = parseMerchants()
	if err != nil {
		return fmt.Errorf("parseMerchants: %w", err)
	}
	err = parseNpcs()
	if err != nil {
		return fmt.Errorf("parseNpcs: %w", err)
	}

	err = parseQuests()
	if err != nil {
		return fmt.Errorf("parseQuests: %w", err)
	}

	err = parseRecipes()
	if err != nil {
		return fmt.Errorf("parseRecipes: %w", err)
	}

	for _, item := range items {
		for _, merchantName := range item.SoldBy {
			if merchantName == "none" || merchantName == "unknown" {
				continue
			}
			merchant := merchants[merchantName]
			if merchant == nil {
				return fmt.Errorf("merchant %s not found", merchantName)
			}
			merchant.Items = append(merchant.Items, item.ID)
		}

		for _, questName := range item.QuestRewarded {
			if questName == "none" || questName == "unknown" {
				continue
			}
			quest := quests[questName]
			if quest == nil {
				return fmt.Errorf("quest %s not found", questName)
			}
		}

		for _, questName := range item.QuestReagent {
			if questName == "none" || questName == "unknown" {
				continue
			}
			quest := quests[questName]
			if quest == nil {
				return fmt.Errorf("quest %s not found", questName)
			}
		}

		for _, npcName := range item.DroppedBy {
			if npcName == "none" || npcName == "unknown" {
				continue
			}
			npc := npcs[npcName]
			if npc == nil {
				return fmt.Errorf("npc %s not found", npcName)
			}
		}

	}

	for _, merchant := range merchants {
		fmt.Printf("%+v\n", merchant)
	}

	for _, quest := range quests {
		fmt.Printf("%+v\n", quest)
	}

	for _, npc := range npcs {
		fmt.Printf("%+v\n", npc)
	}

	return nil
}

func parseMerchants() error {
	err := filepath.WalkDir("db/merchant/", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}
		r, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer r.Close()
		type merchantStruct struct {
			Entry []db.Merchant `yaml:"merchants"`
		}
		merchantFile := &merchantStruct{}
		err = yaml.NewDecoder(r).Decode(merchantFile)
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		for i := range merchantFile.Entry {
			merchant := merchantFile.Entry[i]
			merchant.Name = zonify(merchant.Name, path)
			merchants[merchant.Name] = &merchant
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}
	return nil
}

func parseNpcs() error {
	err := filepath.WalkDir("db/npc/", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}
		r, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer r.Close()
		type npcStruct struct {
			Entry []db.Npc `yaml:"npcs"`
		}
		npcFile := &npcStruct{}
		err = yaml.NewDecoder(r).Decode(npcFile)
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		for i := range npcFile.Entry {
			npc := npcFile.Entry[i]
			npc.Name = zonify(npc.Name, path)
			npcs[npc.Name] = &npc
			//fmt.Printf("%+v\n", npc)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}
	return nil
}

func parseQuests() error {
	err := filepath.WalkDir("db/quest/", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}
		r, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer r.Close()
		type questStruct struct {
			Entry []db.Quest `yaml:"quests"`
		}
		questFile := &questStruct{}
		err = yaml.NewDecoder(r).Decode(questFile)
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		for i := range questFile.Entry {
			quest := questFile.Entry[i]
			quest.ID = zonify(quest.ID, path)
			quests[quest.ID] = &quest
			//fmt.Printf("%+v\n", quest)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}
	return nil
}

func zonify(name string, path string) string {
	zone := strings.TrimSuffix(filepath.Base(path), ".yaml")
	return strings.ReplaceAll(name, " ", "_") + "#" + zone
}

func parseRecipes() error {
	r, err := os.Open("db/recipe.yaml")
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer r.Close()
	type recipeStruct struct {
		Entry []db.Recipe `yaml:"recipes"`
	}
	recipeFile := &recipeStruct{}
	err = yaml.NewDecoder(r).Decode(recipeFile)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	for i := range recipeFile.Entry {
		recipe := recipeFile.Entry[i]
		recipes[recipe.ID] = &recipe
		//fmt.Printf("%+v\n", recipe)
	}
	return nil
}
