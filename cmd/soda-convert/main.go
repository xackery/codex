package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xackery/codex/db"
	"gopkg.in/yaml.v3"
)

type recipeStruct struct {
	Entry []db.Recipe `yaml:"recipes"`
}

var (
	recipeFile = &recipeStruct{}
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed to run:", err)
		os.Exit(1)
	}
}

func run() error {
	err := convertRecipe()
	if err != nil {
		return fmt.Errorf("convert recipe: %w", err)
	}

	err = convertRecipeComponents()
	if err != nil {
		return fmt.Errorf("convert recipe components: %w", err)
	}
	w, err := os.Create("db/recipe.yaml")
	if err != nil {
		return fmt.Errorf("create recipe: %w", err)
	}
	defer w.Close()

	err = yaml.NewEncoder(w).Encode(recipeFile)
	if err != nil {
		return fmt.Errorf("encode recipe: %w", err)
	}
	fmt.Println("Converted", len(recipeFile.Entry), "recipes")

	err = convertItems()
	if err != nil {
		return fmt.Errorf("convert items: %w", err)
	}

	return nil
}

func convertRecipe() error {
	data, err := os.ReadFile("soda/recipes.txt")
	if err != nil {
		return fmt.Errorf("read recipes: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for lineNumber, line := range lines {
		if line == "" {
			continue
		}
		if lineNumber == 0 {
			continue
		}
		recipe := db.Recipe{}
		records := strings.Split(line, "|")
		if len(records) != 13 {
			return fmt.Errorf("line %d: expected 13 records, got %d", lineNumber, len(records))
		}
		recipe.ID, err = strconv.Atoi(records[0])
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		recipe.Name = records[1]
		recipe.Trivial, err = strconv.Atoi(records[2])
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		recipe.RecipeType = records[4]
		recipe.OtherID, err = strconv.Atoi(records[5])
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		recipe.ResultID, err = strconv.Atoi(records[6])
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		recipeFile.Entry = append(recipeFile.Entry, recipe)
	}

	return nil
}

func convertRecipeComponents() error {
	data, err := os.ReadFile("soda/recipecomp.txt")
	if err != nil {
		return fmt.Errorf("read recipe components: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for lineNumber, line := range lines {
		if line == "" {
			continue
		}
		if lineNumber == 0 {
			continue
		}
		recipeComponent := db.RecipeComponent{}
		records := strings.Split(line, "|")
		if len(records) != 12 {
			return fmt.Errorf("line %d: expected 12 records, got %d", lineNumber, len(records))
		}
		recipeComponentID, err := strconv.Atoi(records[0])
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		var recipe *db.Recipe
		for i := range recipeFile.Entry {
			if recipeFile.Entry[i].ID != recipeComponentID {
				continue
			}
			recipe = &recipeFile.Entry[i]
			break
		}
		if recipe == nil {
			return fmt.Errorf("line %d: recipe %d not found", lineNumber, recipeComponentID)
		}

		recipeComponent.ComponentID, err = strconv.Atoi(records[1])
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		recipeComponent.ItemID, err = strconv.Atoi(records[3])
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		recipeComponent.OtherID, err = strconv.Atoi(records[4])
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		recipeComponent.OtherType, err = strconv.Atoi(records[5])
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		recipe.Components = append(recipe.Components, recipeComponent)
	}

	return nil
}

func convertItems() error {
	data, err := os.ReadFile("soda/items.txt")
	if err != nil {
		return fmt.Errorf("read items: %w", err)
	}

	lastRange := 0

	type itemStruct struct {
		Entry []*db.Item `yaml:"items"`
	}
	itemFile := &itemStruct{}

	items := make(map[int]*db.Item)
	err = db.ParseItems(items)
	if err != nil {
		return fmt.Errorf("parse items: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for lineNumber, line := range lines {
		if line == "" {
			continue
		}
		if lineNumber == 0 {
			continue
		}

		item := &db.Item{
			ID:             0,
			Name:           "unknown",
			SoldBy:         []string{"unknown"},
			DroppedBy:      []string{"unknown"},
			QuestRewarded:  []string{"unknown"},
			QuestReagent:   []string{"unknown"},
			RecipeRewarded: []string{"unknown"},
			RecipeReagent:  []string{"unknown"},
		}

		records := strings.Split(line, "|")
		if len(records) < 315 {
			return fmt.Errorf("line %d: expected 315 records, got %d", lineNumber, len(records))
		}
		item.ID, err = strconv.Atoi(records[5])
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		tmpItem, ok := items[item.ID]
		if ok {
			item = tmpItem
		}

		item.Name = records[1]
		if item.ID-lastRange-5000 < 0 {
			itemFile.Entry = append(itemFile.Entry, item)
			continue
		}
		maxValue := item.ID
		// round maxValue down to nearest 1000
		maxValue = maxValue - (maxValue % 1000) - 1

		w, err := os.Create(fmt.Sprintf("db/item/%d-%d.yaml", lastRange, maxValue))
		if err != nil {
			return fmt.Errorf("create item: %w", err)
		}
		defer w.Close()

		enc := yaml.NewEncoder(w)
		enc.SetIndent(2)
		err = enc.Encode(itemFile)
		if err != nil {
			return fmt.Errorf("encode item: %w", err)
		}
		itemFile = &itemStruct{}
		lastRange = item.ID
		itemFile.Entry = append(itemFile.Entry, item)

	}

	return nil
}
