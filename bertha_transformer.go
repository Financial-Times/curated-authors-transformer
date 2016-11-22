package main

import (
	"github.com/jaytaylor/html2text"
	"github.com/pborman/uuid"
)

const tmeAuthority = "http://api.ft.com/system/FT-TME"

type berthaTransformer struct {
}

func (bt *berthaTransformer) authorToPerson(a author) (person, error) {
	uuid := uuid.NewMD5(uuid.UUID{}, []byte(a.TmeIdentifier)).String()
	plainDescription, err := html2text.FromString(a.Biography)

	if err != nil {
		return person{}, err
	}

	altIds := alternativeIdentifiers{
		UUIDS: []string{uuid},
		TME:   []string{a.TmeIdentifier},
	}

	p := person{
		Uuid:                   uuid,
		Name:                   a.Name,
		PrefLabel:              a.Name,
		EmailAddress:           a.Email,
		TwitterHandle:          a.TwitterHandle,
		FacebookProfile:        a.FacebookProfile,
		LinkedinProfile:        a.LinkedinProfile,
		Description:            plainDescription,
		DescriptionXML:         a.Biography,
		ImageUrl:               a.ImageUrl,
		AlternativeIdentifiers: altIds,
	}

	return p, err
}
