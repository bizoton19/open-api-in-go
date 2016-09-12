package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Recall represents info on cpsc recalls.
type Recall struct {
	RecallID        int    `json:"RecallID"`
	RecallNumber    string `json:"RecallNumber"`
	RecallDate      string `json:"RecallDate"`
	Description     string `json:"Description"`
	URL             string `json:"URL"`
	Title           string `json:"Title"`
	ConsumerContact string `json:"ConsumerContact"`
	LastPublishDate string `json:"LastPublishDate"`
	Products        []struct {
		Name          string `json:"Name"`
		Description   string `json:"Description"`
		Model         string `json:"Model"`
		Type          string `json:"Type"`
		CategoryID    string `json:"CategoryID"`
		NumberOfUnits string `json:"NumberOfUnits"`
	} `json:"Products"`
	Inconjunctions []struct{} `json:"Inconjunctions"`
	Images         []struct {
		URL string `json:"URL"`
	} `json:"Images"`
	Injuries []struct {
		Name string `json:"Name"`
	} `json:"Injuries"`
	Manufacturers []struct {
		Name      string `json:"Name"`
		CompanyID string `json:"CompanyID"`
	} `json:"Manufacturers"`
	ManufacturerCountries []struct {
		Country string `json:"Country"`
	} `json:"ManufacturerCountries"`
	ProductUPCs []struct {
		UPC string `json:"UPC"`
	} `json:"ProductUPCs"`
	Hazards []struct {
		Name         string `json:"Name"`
		HazardTypeID string `json:"HazardTypeID"`
	} `json:"Hazards"`
	Remedies []struct {
		Name string `json:"Name"`
	} `json:"Remedies"`
	Retailers []struct {
		Name      string `json:"Name"`
		CompanyID string `json:"CompanyID"`
	} `json:"Retailers"`
}

//GetURL will get the formatted url
func GetURL(recallID string) string {
	var safeURL = url.QueryEscape(recallID)
	var recallURL = "https://www.saferproducts.gov/restwebservices/recall?recallId=" + safeURL + "&format=json"
	return recallURL
}
func main() {
	session, err := mgo.Dial("bilomaxmongo.cloudapp.net:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("CPSC").C("Recalls")
	for recallID := 22900; recallID < 200000; recallID++ {
		var id = fmt.Sprintf("%d", recallID)
		// Build the request
		req, err := http.NewRequest("GET", GetURL(id), nil)
		if err != nil {
			log.Fatal("NewRequest: ", err)
			return
		}
		log.Println(GetURL(id))
		// For control over HTTP client headers,
		// redirect policy, and other settings,
		// create a Client
		// A Client is an HTTP client
		client := &http.Client{}

		// Send the request via a client
		// Do sends an HTTP request and
		// returns an HTTP response
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
			return
		}
		log.Println(resp.ContentLength)

		// Callers should close resp.Body
		// when done reading from it
		// Defer the closing of the body
		defer resp.Body.Close()

		//use json.Decode for reading streams of JSON data
		var unmarshalled = json.NewDecoder(resp.Body)
		//read open bracket
		t, err := unmarshalled.Token()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%T: %v\n", t, t)
		//while the array contains values
		for unmarshalled.More() {
			var recallValue Recall
			err := unmarshalled.Decode(&recallValue)

			if err != nil {
				log.Fatal(err)
				log.Print(resp.Status)

			}
			fmt.Printf("\n%v:\n %v\n%v\n",
				recallValue.RecallID,
				recallValue.Title,
				recallValue.URL)

			

			err = c.Insert(&recallValue)

			if err != nil {
				log.Fatal(err)
			}

			//read closing braket
			t, err = unmarshalled.Token()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%T: %v\n", t, t)

		}

	}

	result := Recall{}
	err = c.Find(bson.M{"recallid": "10"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RecallId:", result.RecallID)

}
