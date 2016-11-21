package main

// This struct reflects the JSON data model of curated authors from Bertha
type author struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	ImageUrl        string `json:"imageurl"`
	Biography       string `json:"biography"`
	TwitterHandle   string `json:"twitterhandle"`
	FacebookProfile string `json:"facebookprofile"`
	LinkedinProfile string `json:"linkedinprofile"`
	TmeIdentifier   string `json:"tmeidentifier"`
}
