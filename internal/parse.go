package internal

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Artists []struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relation     Relation
}

type Relation struct {
	Relations map[string][]string `json:"datesLocations"`
}

type Relations struct {
	Index []Relation `json:"index"`
}

// type DatesLocations struct {
// 	Locations []string `json:"locations"`
// 	Dates     []string `json:"dates"`
// }

// type Dates struct {
// 	ID    int      `json:"id"`
// 	Dates []string `json:"dates"`
// }

// type Locations struct {
// 	ID        int      `json:"id"`
// 	Locations []string `json:"locations"`
// 	Dates     Dates    `json:"dates"`
// }

func ParseJson(w http.ResponseWriter, r *http.Request) (Artists, error) {
	relationData, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Println("Cannot get from URL", err)
		HandlerError(w, 500)
		return nil, err
	}

	defer relationData.Body.Close()

	data, err := ioutil.ReadAll(relationData.Body)
	if err != nil {
		log.Println("Error reading json data:", err)
		return nil, err
	}

	var relations Relations
	err = json.Unmarshal(data, &relations)
	if err != nil {
		HandlerError(w, 500)
		return nil, err
	}

	artistsData, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		HandlerError(w, 500)
		return nil, err
	}

	defer artistsData.Body.Close()
	
	data, err = io.ReadAll(artistsData.Body)
	if err != nil {
		log.Println("Error reading json data:", err)
		return nil, err
	}

	var artists Artists
	err = json.Unmarshal(data, &artists)
	if err != nil {
		HandlerError(w, 500)
		return nil, err
	}

	for i := 0; i < len(artists); i++ {
		artists[i].Relation = relations.Index[i]
	}
	return artists, nil
}
