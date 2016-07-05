package main

type person struct {
	Uuid           string       `json:"uuid"`
	BirthYear      int          `json:"birthYear,omitempty"`
	Identifiers    []identifier `json:"identifiers,omitempty"`
	Name           string       `json:"name,omitempty"`
	Salutation     string       `json:"salutation,omitempty"`
	Aliases        []string     `json:"aliases,omitempty"`
	EmailAddress   string       `json:"emailAddress,omitempty"`
	TwitterHandle  string       `json:"twitterHandle,omitempty"`
	Description    string       `json:"description,omitempty"`
	DescriptionXML string       `json:"descriptionXML,omitempty"`
	ImageUrl       string       `json:"_imageUrl,omitempty"` // TODO this is a temporary thing - needs to be integrated into images properly
}

type identifier struct {
	Authority       string `json:"authority"`
	IdentifierValue string `json:"identifierValue"`
}
