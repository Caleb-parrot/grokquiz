package quiz

// Categories is the start menu. Each query is sent live to Grokipedia;
// the names are not a stored fact pack.
var Categories = []Category{
	{Name: "Space", Queries: []string{
		"Moon", "Mars", "Jupiter", "Saturn", "Venus", "Neptune", "Uranus",
		"Pluto", "Sun", "Earth", "Milky Way", "Black hole", "Solar System",
		"Apollo 11", "Hubble Space Telescope", "International Space Station",
	}},
	{Name: "History", Queries: []string{
		"World War II", "World War I", "Roman Empire", "Ancient Egypt",
		"French Revolution", "Industrial Revolution", "Cold War", "Renaissance",
		"Ottoman Empire", "Mongol Empire", "American Civil War", "Silk Road",
		"Magna Carta", "Byzantine Empire",
	}},
	{Name: "Nature", Queries: []string{
		"Blue whale", "African elephant", "Giant panda", "Cheetah",
		"Emperor penguin", "Great white shark", "Honey bee", "Bald eagle",
		"Amazon rainforest", "Sahara", "Great Barrier Reef", "Nile",
		"Mount Everest", "Komodo dragon", "Grand Canyon",
	}},
	{Name: "Music", Queries: []string{
		"Ludwig van Beethoven", "Wolfgang Amadeus Mozart", "Johann Sebastian Bach",
		"The Beatles", "Jazz", "Opera", "Violin", "Piano", "Symphony", "Guitar", "Blues",
	}},
	{Name: "Sports", Queries: []string{
		"Association football", "Basketball", "Olympic Games", "Tennis",
		"Cricket", "Baseball", "Marathon", "Tour de France", "FIFA World Cup",
		"Ice hockey",
	}},
	{Name: "Places", Queries: []string{
		"Tokyo", "Paris", "Cairo", "Rome", "London", "New York City",
		"Amazon River", "Pacific Ocean", "Antarctica", "Himalayas",
		"Mediterranean Sea", "Australia",
	}},
	{Name: "Science", Queries: []string{
		"Oxygen", "Water", "DNA", "Atom", "Photosynthesis", "Plate tectonics",
		"Periodic table", "Gravity", "Evolution", "Vaccine", "Electricity",
		"Magnetism", "Penicillin",
	}},
	{Name: "Myths", Queries: []string{
		"Zeus", "Hera", "Poseidon", "Hades", "Athena", "Thor", "Odin",
		"Loki", "Heracles", "Ra", "Anubis",
	}},
}
