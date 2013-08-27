package main

import (
	"encoding/xml"
	"fmt"
)

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

func main() {

	var v CharactersResult
	data := `
		<?xml version='1.0' encoding='UTF-8'?>
		<eveapi version="2">
		  <currentTime>2013-08-27 12:43:47</currentTime>
		  <result>
		    <rowset name="characters" key="characterID" columns="name,characterID,corporationName,corporationID">
		      <row name="Daniel Derpington" characterID="93773652" corporationName="Science and Trade Institute" corporationID="1000045" />
		    </rowset>
		  </result>
		  <cachedUntil>2013-08-27 13:40:34</cachedUntil>
		</eveapi>
		`
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Printf("XMLName: %#v\n", v.XMLName)
	fmt.Printf("CurrentTime: %q\n", v.CurrentTime)
	fmt.Printf("CachedUntil: %q\n", v.CachedUntil)
	fmt.Printf("Name: %q\n", v.Characters[0].Name)
	fmt.Printf("ID: %q\n", v.Characters[0].ID)
	fmt.Printf("Corporation: %q\n", v.Characters[0].Corporation)
	fmt.Printf("CorporationID: %q\n", v.Characters[0].CorporationID)
}
