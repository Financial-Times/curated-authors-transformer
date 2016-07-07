package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gregjones/httpcache"
	"net/http"
)

var client = httpcache.NewMemoryCacheTransport().Client()

type berthaService struct {
	berthaUrl  string
	authorsMap map[string]author
}

func (bs *berthaService) getAuthorsCount() (int, error) {
	bs.authorsMap = make(map[string]author)

	resp, err := bs.callBerthaService()
	if err != nil {
		log.Error(err)
		return -1, err
	}

	var authors []author
	if err = json.NewDecoder(resp.Body).Decode(&authors); err != nil {
		log.Error(err)
		return -1, err
	}

	for _, a := range authors {
		bs.authorsMap[a.Uuid] = a
	}
	return len(bs.authorsMap), nil
}

func (bs *berthaService) getAuthorsUuids() []string {
	uuids := make([]string, 0)
	for uuid, _ := range bs.authorsMap {
		uuids = append(uuids, uuid)
	}
	return uuids
}

func (bs *berthaService) getAuthorByUuid(uuid string) author {
	return bs.authorsMap[uuid]
}

func (bs *berthaService) callBerthaService() (res *http.Response, err error) {
	log.WithFields(log.Fields{"bertha_url": bs.berthaUrl}).Info("Calling Bertha...")
	res, err = client.Get(bs.berthaUrl)
	return
}

func (bs *berthaService) checkConnectivity() error {
	resp, err := bs.callBerthaService()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Bertha returns unexpected HTTP status: %d", resp.StatusCode))
	}
	return nil
}
