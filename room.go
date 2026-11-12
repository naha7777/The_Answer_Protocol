package main

type room struct {
	image	string
	player	map[string]string
	pnj		pnj
	north	string
	south	string
	east	string
	west	string
}

type pnj struct {
	image	string
	text	string
}


func rooms() map[string]room {
	roomsmap := map[string]room {
		"map": {
			image: "pictures/Village_Map.png",
			player: nil,
			pnj: pnj{},
			north: "",
			south: "",
			east: "",
			west: "",
		},
		"square": {
			image: "pictures/square.jpeg",
			player: nil,
			pnj: pnj{
				image: "pictures/pnj/poulet.png",
				text: "cot cot cot",
			},
			north: "tavern",
			south: "",
			east: "market",
			west: "",
		},
		"tavern": {
			image: "pictures/tavern.jpeg",
			player: nil,
			pnj: pnj{},
			north: "haunted_forest",
			south: "square",
			east: "flowering_meadow",
			west: "secret_garden",
		},
		"flowering_meadow": {
			image: "pictures/flowering_meadow.jpeg",
			player: nil,
			pnj: pnj{
				image: "pictures/pnj/fee.png",
				text: "blablabla",
			},
			north: "",
			south: "market",
			east: "",
			west: "tavern",
		},
		"market": {
			image: "pictures/market.jpeg",
			player: nil,
			pnj: pnj{},
			north: "flowering_meadow",
			south: "",
			east: "",
			west: "square",
		},
		"secret_garden": {
			image: "pictures/secret_garden.jpeg",
			player: nil,
			pnj: pnj{},
			north: "lake",
			south: "",
			east: "tavern",
			west: "",
		},
		"lake": {
			image: "pictures/lake.jpeg",
			player: nil,
			pnj: pnj{},
			north: "",
			south: "secret_garden",
			east: "haunted_forest",
			west: "",
		},
		"haunted_forest": {
			image: "pictures/haunted_forest.jpeg",
			player: nil,
			pnj: pnj{},
			north: "castle",
			south: "tavern",
			east: "",
			west: "lake",
		},
		"castle": {
			image: "pictures/castle.jpeg",
			player: nil,
			pnj: pnj{},
			north: "",
			south: "haunted_forest",
			east: "",
			west: "",
		},
	}
	return roomsmap
}