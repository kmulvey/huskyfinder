package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type BoulderResponse struct {
	Results BoulderAnimal `json:"adoptableSearch"`
}

type BoulderAnimal struct {
	ID                   string        `json:"ID"`
	Name                 string        `json:"Name"`
	Species              string        `json:"Species"`
	Sex                  string        `json:"Sex"`
	PrimaryBreed         string        `json:"PrimaryBreed"`
	SecondaryBreed       interface{}   `json:"SecondaryBreed"`
	Sn                   interface{}   `json:"SN"`
	Age                  string        `json:"Age"`
	Photo                string        `json:"Photo"`
	Location             string        `json:"Location"`
	OnHold               string        `json:"OnHold"`
	SpecialNeeds         []string      `json:"SpecialNeeds"`
	NoDogs               []string      `json:"NoDogs"`
	NoCats               []string      `json:"NoCats"`
	NoKids               []string      `json:"NoKids"`
	MemoList             []interface{} `json:"MemoList"`
	Arn                  []string      `json:"ARN"`
	BehaviorTestList     []string      `json:"BehaviorTestList"`
	Stage                string        `json:"Stage"`
	AnimalType           string        `json:"AnimalType"`
	AgeGroup             string        `json:"AgeGroup"`
	WildlifeIntakeInjury []string      `json:"WildlifeIntakeInjury"`
	WildlifeIntakeCause  []string      `json:"WildlifeIntakeCause"`
	BuddyID              string        `json:"BuddyID"`
	Featured             string        `json:"Featured"`
	Sublocation          string        `json:"Sublocation"`
	ChipNumber           interface{}   `json:"ChipNumber"`
	FreshnessStamp       string        `json:"FreshnessStamp"`
}

type LongmontResponse struct {
	Success bool             `json:"success"`
	Data    []LongmontAnimal `json:"data"`
}

type LongmontAnimal struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Species      string      `json:"species"`
	Sex          string      `json:"sex"`
	PrimaryBreed string      `json:"primary_breed"`
	Sn           interface{} `json:"sn"`
	Age          string      `json:"age"`
	AgeFormatted string      `json:"age_formatted"`
	AgeGroup     string      `json:"age_group"`
	PhotoURL     string      `json:"photo_url"`
	Location     string      `json:"location"`
	OnHold       string      `json:"on_hold"`
	NoDogs       bool        `json:"no_dogs"`
	NoCats       bool        `json:"no_cats"`
	ChipNumber   interface{} `json:"chip_number"`
	IsBarnCat    bool        `json:"is_barn_cat"`
}

type PetangoResponse struct {
	Items     []PetangoAnimal `json:"items"`
	Count     int             `json:"count"`
	ZipCode   string          `json:"zipCode"`
	BreedName string          `json:"breedName"`
}

type PetangoAnimal struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	SpeciesID    int         `json:"speciesId"`
	Photo        string      `json:"photo"`
	Age          string      `json:"age"`
	Distance     int         `json:"distance"`
	Gender       string      `json:"gender"`
	Breed        string      `json:"breed"`
	NoDogs       bool        `json:"noDogs"`
	NoCats       bool        `json:"noCats"`
	NoKids       bool        `json:"noKids"`
	URL          string      `json:"url"`
	AdoptionDate interface{} `json:"adoptionDate"`
	HasVideo     bool        `json:"hasVideo"`
	Score        int         `json:"score"`
}

func main() {
	boulder()
	longmont()
	petango()
}

func boulder() {
	res, err := http.Get("https://boulderhumane.org/wp-content/plugins/Petpoint-Webservices-2018/pullanimals.php?type=dog")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	var dogs []BoulderResponse
	err = json.NewDecoder(res.Body).Decode(&dogs)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range dogs {
		// Secondary breed is not consistantly a string for some reason
		var secondary, _ = d.Results.SecondaryBreed.(string) // nothing we can do about the error
		if strings.Contains(strings.ToLower(d.Results.PrimaryBreed), "husky") || strings.Contains(strings.ToLower(secondary), "husky") {
			fmt.Println("Boulder", d.Results.Name)
		}
	}

	longmont()
}

func longmont() {
	endpoint := "https://www.longmonthumane.org/wp-admin/admin-ajax.php"
	data := url.Values{}
	data.Set("action", "search_adoptable")

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var results LongmontResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range results.Data {
		if strings.Contains(strings.ToLower(d.Species), "dog") {
			if strings.Contains(strings.ToLower(d.PrimaryBreed), "husky") {
				fmt.Println("Longmont", d.Name)
			}
		}
	}
}

func petango() {
	endpoint := "https://www.petango.com/DesktopModules/Pethealth.Petango/Pethealth.Petango.DnnModules.AnimalSearchResult/API/Main/Search"
	data := url.Values{}
	data.Set("action", "search_adoptable")

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(`{"location":"80301","distance":"200","speciesId":"1","breedId":"670","goodWithDogs":false,"goodWithCats":false,"goodWithChildren":false,"mustHavePhoto":        false,"mustHaveVideo":false,"happyTails":false,"lostAnimals":false,"moduleId":843,"recordOffset":0,"recordAmount":26}`)))
	request.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("ModuleId", "843")
	request.Header.Add("TabId", "260")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var results PetangoResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range results.Items {
		if strings.Contains(strings.ToLower(d.Breed), "husky") { // i think the req body takes care of this but cant hurt anyway
			fmt.Println("Petango", d.Name)
		}
	}
}
