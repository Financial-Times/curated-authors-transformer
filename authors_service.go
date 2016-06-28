package main

type authorsService interface {
	getAuthorsUuids() ([]string, error)
	getAuthorByUuid(uuid string) author
	checkConnectivity() error
}
