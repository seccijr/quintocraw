package pisoscom

import (
	"io"
	"log"
	"github.com/seccijr/quintocrawl/model"
	"github.com/seccijr/quintocrawl/schemas/pisoscom/page"
	"github.com/seccijr/quintocrawl/client"
)

type PCBroker struct {
	flats    model.FlatRepo
	url string
}

func (broker *PCBroker) Parse(httpBody io.Reader) []*client.Broker {
	var links []string
	var brokers []client.Broker
	pcg, err := page.NewDocFromReader(httpBody)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	switch {
	case pcg.HasStates():
		links, _ = pcg.GetStateLinks()
		break
	case pcg.HasProvinces():
		links, _ = pcg.GetProvinceLinks()
		break
	}

	for link := range links {
		brokers = append(&PCBroker{flats: broker.flats, url: link})
	}

	return brokers, nil
}

func (broker *PCBroker) URL() string {
	return broker.url
}
