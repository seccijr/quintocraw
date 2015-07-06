package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"github.com/seccijr/quintocrawl/model/mongo"
	"github.com/seccijr/quintocrawl/schemas/pisoscom"
	"github.com/seccijr/quintocrawl/client"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("flat")
	flatsRepo := mongo.MFlatRepo{c}
	pcConfig, err := pisoscom.ReadConfig("schemas/pisoscom/params.json")

	if err != nil {
		fmt.Println("Error reading Pisos.com config")
	}

	pcUrl := pcConfig.Base.String()
	pcBroker := &pisoscom.PCBroker{flatsRepo, pcUrl, "/"}
	client := client.New()
	client.Handle(pcBroker)
	client.Run()
}
