package main

type transformer interface {
	authorToPerson(author) (person, error)
}
