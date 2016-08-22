package main

type authorsService interface {
	refreshCache() error
	getAuthorsCount() int
	getAuthorsUuids() []string
	getAuthorByUuid(uuid string) person
	checkConnectivity() error
}
