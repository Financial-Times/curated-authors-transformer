package main

import (
	"github.com/jaytaylor/html2text"
)

const tmeAuthority = "http://api.ft.com/system/FT-TME"

type berthaTransformer struct {
}

func (bt *berthaTransformer) authorToPerson(a author) (person, error) {
	plainDescription, err := html2text.FromString(a.Biography)

	if err != nil {
		return person{}, err
	}

	id := identifier{
		Authority:       tmeAuthority,
		IdentifierValue: a.TmeIdentifier,
	}

	identifiers := []identifier{id}

	p := person{
		Uuid:           a.Uuid,
		Name:           a.Name,
		EmailAddress:   a.Email,
		TwitterHandle:  a.TwitterHandle,
		Description:    plainDescription,
		DescriptionXML: a.Biography,
		ImageUrl:       a.ImageUrl,
		Identifiers:    identifiers,
	}

	return p, err
}
