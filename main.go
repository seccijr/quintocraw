package main

import (
	"github.com/seccijr/quintocrawl/client"
	"github.com/seccijr/quintocrawl/schemas/pisoscom"
	"github.com/seccijr/quintocrawl/model/mongo"
	"fmt"
)

func main() {
	flatsRepo := mongo.MFlatRepo{}
	pcConfig, err := pisoscom.ReadConfig("schemas/pisoscom/params.json")

	if err != nil {
		fmt.Println("Error reading Pisos.com config")
		return -1
	}

	pcUrl := pcConfig.Base.String()
	pcBroker := pisoscom.PCBroker{flats: flatsRepo, url: pcUrl}
	client := client.New()
	client.Handle(pcBroker)

	return 0
}
