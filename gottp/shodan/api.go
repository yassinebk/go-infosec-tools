package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIInfo struct {
	QueryCredits int    `json:"query_credits"`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet"`
	Plan         string `json:"plan"`
	Https        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
}

func (s *Client) APIInfo() (*APIInfo, error) {
	res, err := http.Get(fmt.Sprintf("%s/api-info?key=%s", BASE_URL, s.apiKey))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var apiInfo APIInfo
	if err := json.NewDecoder(res.Body).Decode(&apiInfo); err != nil {
		return nil, err
	}

	return &apiInfo, nil
}

func (s *Client) HostSearch(q string) (*HostSearch, error) {
	res, err := http.Get(fmt.Sprintf("%s/shodan/host/search?key=%s&query=%s", BASE_URL, s.apiKey, q))

	if err != nil {
		return nil, err
		// log.Fatalln("Error querying for data", err)
	}

	defer res.Body.Close()

	var hostSearch HostSearch
	if err := json.NewDecoder(res.Body).Decode(&hostSearch); err != nil {
		return nil, err
	}

	return &hostSearch, nil
}
