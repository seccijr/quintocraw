package pisoscom

import (
	"io"
	"log"
	"github.com/seccijr/quintocrawl/model"
	"github.com/seccijr/quintocrawl/schemas/pisoscom/page"
	"github.com/seccijr/quintocrawl/client"
)

type PCBroker struct {
	Flats model.FlatRepo
	Base  string
	Url   string
}

func (broker PCBroker) Parse(httpBody io.Reader) ([]client.Broker, error) {
	var links []string
	var brokers []client.Broker
	pcg, err := page.NewDocFromReader(httpBody)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	wStates, _ := pcg.HasStates()
	wProvinces, _ := pcg.HasProvinces()

	switch {
	case wStates:
		links, _ = pcg.GetStateLinks()
		break
	case wProvinces:
		links, _ = pcg.GetProvinceLinks()
		break
	}

	for _, link := range links {
		newBroker := PCBroker{broker.Flats, broker.Base, link}
		brokers = append(brokers, &newBroker)
	}

	return brokers, nil
}

func (broker PCBroker) URL() string {
	newUrl := client.FixUrl(broker.Url, broker.Base)
	return newUrl
}
