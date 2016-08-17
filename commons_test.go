package main

var martinWolfUuid = "0f07d468-fc37-3c44-bf19-a81f2aae9f36"
var lucyKellawayUuid = "8f9ac45f-2cc2-35f7-83f4-579c66a09eb0"

var uuids = []string{martinWolfUuid, lucyKellawayUuid}

var martinWolf = author{
	Name:          "Martin Wolf",
	Email:         "martin.wolf@ft.com",
	ImageUrl:      "https://next-geebee.ft.com/image/v1/images/raw/fthead:martin-wolf?source=next",
	Biography:     "Martin Wolf is chief economics commentator at the Financial Times, London.",
	TwitterHandle: "@martinwolf_",
	TmeIdentifier: "Q0ItMDAwMDkwMA==-QXV0aG9ycw==",
}

var lucyKellaway = author{
	Name:          "Lucy Kellaway",
	Email:         "lucy.kellaway@ft.com",
	ImageUrl:      "https://next-geebee.ft.com/image/v1/images/raw/fthead:lucy-kellaway?source=next",
	Biography:     "Lucy Kellaway is an Associate Editor and management columnist of the FT. For the past 15 years her weekly Monday column has poked fun at management fads and jargon and celebrated the ups and downs of office life.",
	TmeIdentifier: "Q0ItMDAwMDkyNg==-QXV0aG9ycw==",
}

var transformedMartinWolf = person{
	Uuid:                   martinWolfUuid,
	Name:                   "Martin Wolf",
	PrefLabel:              "Martin Wolf",
	EmailAddress:           "martin.wolf@ft.com",
	TwitterHandle:          "@martinwolf_",
	Description:            "Martin Wolf is chief economics commentator at the Financial Times, London.",
	DescriptionXML:         `<p>Martin Wolf is chief economics commentator at the Financial Times, London.</p>`,
	ImageUrl:               "https://next-geebee.ft.com/image/v1/images/raw/fthead:martin-wolf?source=next",
	AlternativeIdentifiers: martinWolfAltIds,
}

var martinWolfAltIds = alternativeIdentifiers{
	TME:   []string{martinWolf.TmeIdentifier},
	UUIDS: []string{martinWolfUuid},
}
