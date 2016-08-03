package main

type authorsService interface {
	getAuthorsCount() (int, error)
	getAuthorsUuids() []string
	getAuthorByUuid(uuid string) person
	checkConnectivity() error
}
