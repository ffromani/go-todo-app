package uuid

import (
	"encoding/json"
	"io"
	"net/http"
)

type UUIDGenerator struct {
	url string
}

var uuidURL = "https://www.uuidtools.com/api/generate/v1/"

func New() UUIDGenerator {
	return UUIDGenerator{url: uuidURL}
}

func (g UUIDGenerator) NewUUID() (string, error) {
	resp, err := http.Get(g.url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var uuids []string
	if err := json.Unmarshal(body, &uuids); err != nil {
		return "", err
	}
	return uuids[0], nil
}
