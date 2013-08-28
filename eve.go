package eve

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)
const eveUrl string = "https://api.eveonline.com"

var charactersCache = make(map[string]CharactersResult)

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
func Characters(keyID string, vCode string) (res []Character) {
	var v CharactersResult

	cacheKey := fmt.Sprintf("%s:%s", keyID, vCode)
	if cachedResult, ok := charactersCache[cacheKey]; ok {
		if cacheEntryValid(cachedResult) {
			fmt.Println("Found cached character result")
			return cachedResult.Characters
		}
	}

	url := fmt.Sprintf("%s/account/Characters.xml.aspx?keyID=%s&vCode=%s", eveUrl, keyID, vCode)
	data, err := fetch(url)
	if err != nil {
		fmt.Printf("Fetch error: %v\n", err)
		return
	}

	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		return
	}

	fmt.Println("Storing result into cache:", cacheKey)
	charactersCache[cacheKey] = v
	return v.Characters
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

func cacheEntryValid(entry CharactersResult) bool {
	t, err := time.Parse("2006-01-02 15:04:05", entry.CachedUntil)
	if err != nil {
		fmt.Println(err)
	}

	if t.After(time.Now()) {
		return true
	}
	return false
}
