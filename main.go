package main

import (
	"gopkg.in/mgo.v2"
	"github.com/seccijr/quintocrawl/model/mongo"
	"github.com/seccijr/quintocrawl/schema/pisoscom"
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
	pcConfig, err := pisoscom.ReadConfig("./schema/pisoscom/params.json")

	if err != nil {
		panic(err)
	}

	pcUrl := pcConfig.Base.String()
	pcBroker := pisoscom.PCBroker{flatsRepo, pcUrl}
	queue := client.NewQueue(&pcBroker)
	queue.Push(pcBroker.Format(""))
	queue.Run()
}
