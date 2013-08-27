package eve

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"io/ioutil"
)
const eveUrl string = "https://api.eveonline.com"

var charactersCache map[string]CharactersResult

type Character struct {
	Name	string	`xml:"name,attr"`
	ID	string	`xml:"characterID,attr"`
	Corporation	string	`xml:"corporationName,attr"`
	CorporationID	string	`xml:"corporationID,attr"`
}

type CharactersResult struct {
	XMLName		xml.Name	`xml:"eveapi"`
	CurrentTime	string		`xml:"currentTime"`
	Characters	[]Character	`xml:"result>rowset>row"`
	CachedUntil	string		`xml:"cachedUntil"`
}

// Characters fetches, parses and returns a set of Eve Online characters.
func Characters(keyID string, vCode string) (res CharactersResult) {
	var v CharactersResult

	cacheKey := fmt.Sprintf("%s:%s", keyID, vCode)
	if cachedResult, ok := charactersCache[cacheKey]; ok {
		if cacheEntryValid(cachedResult) {
			fmt.Printf("Found cached character result")
			return cachedResult
		}
	}

	url := fmt.Sprintf("%s/account/Characters.xml.aspx?keyID=%s&vCode=%s", eveUrl, keyID, vCode)
	data, err := fetch(url)
	if err != nil {
		fmt.Printf("Fetch error: %v", err)
		return
	}

	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("Unmarshal error: %v", err)
		return
	}

	charactersCache[cacheKey] = v
	return v
}

func fetch(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}
