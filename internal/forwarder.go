package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// pushes the metadata for server discovery to the replicate server
func SendDiscover(cfg Config, md *Metadata) error {
	body, err := json.Marshal(md)
	if err != nil {
		return err
	}

	url := cfg.ServerUrl + "/discover"
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("discover endpoint returned %s", res.Status)
	}
	return nil
}

