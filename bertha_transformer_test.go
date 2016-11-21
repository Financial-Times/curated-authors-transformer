package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var aBioXml = `<p>Eric Theodore Cartman is one of the main characters in the animated television series <a href="https://en.wikipedia.org/wiki/South_Park">South Park</a>, created by Matt Stone and Trey Parker, and voiced by Trey Parker.</p>`
var aBio = "Eric Theodore Cartman is one of the main characters in the animated television series South Park ( https://en.wikipedia.org/wiki/South_Park ) , created by Matt Stone and Trey Parker, and voiced by Trey Parker."

var cartmanUuid = "bbd8c19f-8f7d-33ae-8b0a-ad65f03e951a"

var anAuthor = author{
	Name:            "Eric Cartman",
	Email:           "eric.cartman@southpark.cc.com",
	ImageUrl:        "https://upload.wikimedia.org/wikipedia/en/7/77/EricCartman.png",
	Biography:       aBioXml,
	TwitterHandle:   "@SouthPark",
	FacebookProfile: "OfficialCartman",
	LinkedinProfile: "ProfessionalCartman",
	TmeIdentifier:   "Q0ItMDAwMDkwMA==-QXV0aG8ycw==",
}

var someAltIds = alternativeIdentifiers{
	TME:   []string{anAuthor.TmeIdentifier},
	UUIDS: []string{cartmanUuid},
}

var aPerson = person{
	Uuid:                   cartmanUuid,
	Name:                   "Eric Cartman",
	PrefLabel:              "Eric Cartman",
	EmailAddress:           "eric.cartman@southpark.cc.com",
	TwitterHandle:          "@SouthPark",
	FacebookProfile:        "OfficialCartman",
	LinkedinProfile:        "ProfessionalCartman",
	Description:            aBio,
	DescriptionXML:         aBioXml,
	ImageUrl:               "https://upload.wikimedia.org/wikipedia/en/7/77/EricCartman.png",
	AlternativeIdentifiers: someAltIds,
}

func TestShouldTransformAuthorToPersonSucessfully(t *testing.T) {
	transformer := berthaTransformer{}
	p, err := transformer.authorToPerson(anAuthor)
	assert.Nil(t, err)
	assert.Equal(t, aPerson, p, "The author")
}
