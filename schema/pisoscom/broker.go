package pisoscom

import (
	"io"
	"log"
	"github.com/seccijr/quintocrawl/model"
	"github.com/seccijr/quintocrawl/schema/pisoscom/page"
	"github.com/seccijr/quintocrawl/client"
)

type PCBroker struct {
	Flats model.FlatRepo
	Base  string
}

func (broker PCBroker) Parse(httpBody io.Reader) ([]string, error) {
	pcg, err := page.NewDocFromReader(httpBody)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return broker.hub(pcg), nil
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
//		var flat model.Flat
//		flat, err = pcg.ParseDetail()
		_, err = pcg.ParseDetail()
		if err == nil {
//			err = broker.Flats.Save(flat)
		}
		if err != nil {
			panic(err)
		}
		break
	}

	return links
}

func (broker PCBroker) Format(url string) string {
	newUrl := client.FixUrl(url, broker.Base)
	return newUrl
}