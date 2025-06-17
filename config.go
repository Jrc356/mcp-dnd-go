package main

const (
	APIBaseURL            = "https://www.dnd5eapi.co/api"
	RequestTimeoutSeconds = 10
)

var CategoryDescriptions = map[string]string{
	"ability-scores":       "The six abilities that describe a character's physical and mental characteristics",
	"alignments":           "The moral and ethical attitudes and behaviors of creatures",
	"backgrounds":          "Character backgrounds and their features",
	"classes":              "Character classes with features, proficiencies, and subclasses",
	"conditions":           "Status conditions that affect creatures",
	"damage-types":         "Types of damage that can be dealt",
	"equipment":            "Items, weapons, armor, and gear for adventuring",
	"equipment-categories": "Categories of equipment",
	"feats":                "Special abilities and features",
	"features":             "Class and racial features",
	"languages":            "Languages spoken throughout the multiverse",
	"magic-items":          "Magical equipment with special properties",
	"magic-schools":        "Schools of magic specialization",
	"monsters":             "Creatures and foes",
	"proficiencies":        "Skills and tools characters can be proficient with",
	"races":                "Character races and their traits",
	"rule-sections":        "Sections of the game rules",
	"rules":                "Game rules",
	"skills":               "Character skills tied to ability scores",
	"spells":               "Magic spells with effects, components, and descriptions",
	"subclasses":           "Specializations within character classes",
	"subraces":             "Variants of character races",
	"traits":               "Racial traits",
	"weapon-properties":    "Special properties of weapons",
}
