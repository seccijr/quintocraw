package main

import (
	"github.com/seccijr/quintocrawl/schemas/pisoscom"
	"github.com/seccijr/quintocrawl/model/mongo"
	"fmt"
)

func main() {
	flatsRepo := mongo.MFlatRepo{}
	pcParser := pisoscom.PCParser{flats: flatsRepo}
	fmt.Print(pcParser)
}
