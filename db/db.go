package db

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
