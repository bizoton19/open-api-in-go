package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	//register a hndler for /recalls/1
	initAPI()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), initAPI()))
}

func initAPI() http.Handler {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/recalls/:id", GetRecall),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	return api.MakeHandler()
}

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

//GetRecall returns a spec recall
func GetRecall(w rest.ResponseWriter, r *rest.Request) {
	recallid, err := strconv.ParseUint(r.PathParam("id"), 10, 32)
	if err != nil {
		rest.NotFound(w, r)
		return
	}

	session, err := mgo.Dial("bilomaxmongo.cloudapp.net:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("CPSC").C("Recalls")

	result := Recall{}
	err = c.Find(bson.M{"recallid": recallid}).One(&result)
	if err != nil {
		rest.Error(w, "Not found in GetRecall function", http.StatusNotFound)
		return
	}

	w.WriteJson(result)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "We Have Recalls")
}
