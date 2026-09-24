package quiz

// Categories is the start menu. Each question searches one phrase from
// that category, at a random offset.
var Categories = []Category{
	{Name: "Space", Queries: []string{
		"planet", "galaxy", "nebula", "exoplanet", "comet", "asteroid",
		"black hole", "constellation", "space probe", "moon",
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
	{Name: "War", Queries: []string{
		"war", "battle", "siege", "world war", "civil war", "naval battle",
	}},
	{Name: "Wildlife", Queries: []string{
		"mammal", "bird of prey", "reptile", "amphibian", "insect", "shark",
		"endangered species", "whale",
	}},
	{Name: "Earth", Queries: []string{
		"volcano", "glacier", "rainforest", "desert", "earthquake", "tornado",
		"canyon", "waterfall",
	}},
	{Name: "Music", Queries: []string{
		"composer", "jazz", "opera", "orchestra", "classical guitar", "symphony",
		"musical instrument",
	}},
	{Name: "Sports", Queries: []string{
		"Olympic Games", "association football", "basketball", "tennis",
		"marathon", "cricket", "World Cup",
	}},
	{Name: "Places", Queries: []string{
		"capital city", "landmark", "island", "castle", "cathedral", "ruins",
		"world heritage",
	}},
	{Name: "Science", Queries: []string{
		"physics", "chemical element", "genetics", "human anatomy",
		"invention", "scientist",
	}},
	{Name: "Myths", Queries: []string{
		"Zeus", "Odin", "dragon", "folklore", "King Arthur",
		"Greek myth", "Norse myth",
	}},
}
