package main

import (
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

func main() {
	boulder()
	longmont()
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
