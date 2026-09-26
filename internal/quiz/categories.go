package quiz

// Categories is the start menu. Each question searches one phrase from
// that category, at a random offset.
var Categories = []Category{
	{Name: "Astronomy & Space", Queries: []string{
		"planet", "galaxy", "nebula", "exoplanet", "comet", "asteroid",
		"black hole", "constellation", "space probe", "moon",
	}},
	{Name: "Earth Science", Queries: []string{
		"volcano", "glacier", "rainforest", "desert", "earthquake", "tornado",
		"canyon", "waterfall",
	}},
	{Name: "Life Science", Queries: []string{
		"biology", "virus", "bacteria", "protist",
		"algae", "fungus", "plant",
		"invertebrate", "insect",
		"mammal", "bird", "fish", "reptile", "amphibian", "shark", "whale",
		"human anatomy", "physiology", "medicine",
		"ecosystem", "endangered species",
		"agriculture", "livestock",
		"biologist",
	}},
	{Name: "Physical Science", Queries: []string{
		"physics", "chemical element", "mathematics", "invention", "scientist",
	}},
	{Name: "History", Queries: []string{
		"Egyptian pharaoh", "Mali Empire", "Carthage",
		"Greek Empire", "Roman Empire", "Viking Age", "French Revolution", "Byzantine Empire",
		"Han dynasty", "Mongol Empire", "Meiji Restoration", "Mughal Empire",
		"American Revolution", "Aztec Empire",
		"Inca Empire", "Simón Bolívar", "Machu Picchu",
		"Māori", "Australian Aboriginal", "Captain Cook",
		"Mesopotamia", "Ottoman Empire", "Persian Empire",
	}},
	{Name: "Geography", Queries: []string{
		"country", "mountain range", "landform", "river", "lake",
		"ocean", "sea", "island",
		"U.S. state", "Canadian province",
		"region of the world", "world city", "capital city",
		"national park", "monument", "landmark",
		"castle", "cathedral", "ruins", "world heritage",
		"cartography", "explorer",
	}},
	{Name: "Performing Arts", Queries: []string{
		"composer", "jazz", "opera", "orchestra", "classical guitar", "symphony",
		"musical instrument",
	}},
	{Name: "Sports", Queries: []string{
		"Olympic Games", "association football", "basketball", "tennis",
		"marathon", "cricket", "World Cup",
	}},
	{Name: "Myths", Queries: []string{
		"Zeus", "Odin", "dragon", "folklore", "King Arthur",
		"Greek myth", "Norse myth",
	}},
}
