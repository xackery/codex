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
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed run: ", err)
		os.Exit(1)
	}
}

func run() error {
	err := parseItems()
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

	for _, item := range items {
		for _, merchantName := range item.SoldBy {
			if merchantName == "none" {
				continue
			}
			merchant := merchants[merchantName]
			if merchant == nil {
				return fmt.Errorf("merchant %s not found", merchantName)
			}
			merchant.Items = append(merchant.Items, item.ID)
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

func parseItems() error {
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
		r, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer r.Close()
		type itemStruct struct {
			Entry []db.Item `yaml:"items"`
		}
		itemFile := &itemStruct{}
		err = yaml.NewDecoder(r).Decode(itemFile)
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		for i := range itemFile.Entry {
			item := itemFile.Entry[i]
			items[item.ID] = &item
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
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
		basePath := strings.TrimSuffix(filepath.Base(path), ".yaml")
		for i := range merchantFile.Entry {
			merchant := merchantFile.Entry[i]
			merchant.Name = merchant.Name + "#" + basePath
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
