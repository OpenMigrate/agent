package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscoverInfo struct {
  Id string `json:"id"`
}
// pushes the metadata for server discovery to the replicate server
func SendDiscover(cfg Config, md *Metadata) (DiscoverInfo, error) {
  disco := DiscoverInfo{}
	body, err := json.Marshal(md)
	if err != nil {
		return disco, err
	}

	url := cfg.ServerUrl + "/discover"
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return disco, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return disco, fmt.Errorf("discover endpoint returned %s", res.Status)
	}

  json.NewDecoder(res.Body).Decode(&disco)
	return disco, nil
}

func ValidateDiscovery(cfg Config, uuid string) (err error){
  url := cfg.ServerUrl + "/server/" + uuid
  res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
    return fmt.Errorf("Discovery failed: uuid not found", res.Status)
	}
  return
}
