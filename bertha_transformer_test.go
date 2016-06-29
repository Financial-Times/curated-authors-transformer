package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var aBioXml = `<p>Eric Theodore Cartman is one of the main characters in the animated television series <a href="https://en.wikipedia.org/wiki/South_Park">South Park</a>, created by Matt Stone and Trey Parker, and voiced by Trey Parker.</p>`
var aBio = "Eric Theodore Cartman is one of the main characters in the animated television series South Park ( https://en.wikipedia.org/wiki/South_Park ) , created by Matt Stone and Trey Parker, and voiced by Trey Parker."
var anIdentifier = identifier{tmeAuthority, "Q0ItMDAwMDkwMA==-QXV0aG8ycw=="}

var anAuthor = author{
	Name:          "Eric Cartman",
	Email:         "eric.cartman@southpark.cc.com",
	ImageUrl:      "https://upload.wikimedia.org/wikipedia/en/7/77/EricCartman.png",
	Biography:     aBioXml,
	TwitterHandle: "@SouthPark",
	Uuid:          "4a893fa2-e58b-4c28-aa12-4bb469cd7e57",
	TmeIdentifier: "Q0ItMDAwMDkwMA==-QXV0aG8ycw==",
}

var aPerson = person{
	Uuid:           "4a893fa2-e58b-4c28-aa12-4bb469cd7e57",
	Name:           "Eric Cartman",
	EmailAddress:   "eric.cartman@southpark.cc.com",
	TwitterHandle:  "@SouthPark",
	Description:    aBio,
	DescriptionXML: aBioXml,
	ImageUrl:       "https://upload.wikimedia.org/wikipedia/en/7/77/EricCartman.png",
	Identifiers:    []identifier{anIdentifier},
}

func TestShouldTransformAuthorToPersonSucessfully(t *testing.T) {
	transformer := berthaTransformer{}
	p, err := transformer.authorToPerson(anAuthor)
	assert.Nil(t, err)
	assert.Equal(t, aPerson, p, "The author")
}
