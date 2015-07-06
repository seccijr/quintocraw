package pisoscom

import (
	"io"
	"log"
	"github.com/seccijr/quintocrawl/model"
	"github.com/seccijr/quintocrawl/schemas/pisoscom/page"
	"github.com/seccijr/quintocrawl/client"
	"fmt"
)

type PCBroker struct {
	Flats model.FlatRepo
	Base  string
	Url   string
}

func (broker PCBroker) Parse(httpBody io.Reader) ([]client.Broker, error) {
	var brokers []client.Broker
	pcg, err := page.NewDocFromReader(httpBody)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	links := broker.hub(pcg)
	for _, link := range links {
		newBroker := PCBroker{broker.Flats, broker.Base, link}
		brokers = append(brokers, &newBroker)
	}

	return brokers, nil
}

func (broker PCBroker) hub(pcg *page.PCDoc) []string {
	var links []string

	switch {
	case pcg.HasStates():
		links, _ = pcg.GetStateLinks()
		break
	case pcg.HasProvinces():
		links, _ = pcg.GetProvinceLinks()
		break
	case pcg.HasZones():
		links, _ = pcg.GetZoneLinks()
		break
	case pcg.HasAds():
		links, _ = pcg.GetAdvLinks()
		break
	case pcg.IsDetail():
		var err error
		var flat model.Flat
		flat, err = pcg.ParseDetail()
		if err == nil {
			err = broker.Flats.Save(flat)
		}
		if err != nil {
			fmt.Println(err)
		}
		break
	}

	return links
}

func (broker PCBroker) URL() string {
	newUrl := client.FixUrl(broker.Url, broker.Base)
	return newUrl
}
