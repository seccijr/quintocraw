package page

import (
	"github.com/PuerkitoBio/goquery"
)

const STATE_SELECT = ".zonas.clearfix a.bold"
const PROV_SELECT = ".zonas.clearfix a[class=\"\"]"
const ZONE_SELECT = ".zoneList a"
const ADV_SELECT = ".gridList a.anuncioLink"
const DETAIL_SELECT = ".content.Detalle"

func (doc *PCDoc) HasStates() bool {
	wStates, _ := Has(doc.dom, STATE_SELECT)
	return wStates
}

func (doc *PCDoc) HasProvinces() bool {
	wProvinces, _ := Has(doc.dom, PROV_SELECT)
	return wProvinces
}

func (doc *PCDoc) HasZones() bool {
	wCounties , _ := Has(doc.dom, ZONE_SELECT)
	return wCounties
}

func (doc *PCDoc) HasAds() bool {
	wMainZone, _ := Has(doc.dom, ADV_SELECT)
	return wMainZone
}

func (doc *PCDoc) IsDetail() bool {
	wMainZone, _ := Has(doc.dom, DETAIL_SELECT)
	return wMainZone
}

func (doc *PCDoc) GetStateLinks() ([]string, error) {
	return GetLinks(doc.dom, STATE_SELECT)
}

func (doc *PCDoc) GetProvinceLinks() ([]string, error) {
	return GetLinks(doc.dom, PROV_SELECT)
}

func (doc *PCDoc) GetZoneLinks() ([]string, error) {
	return GetLinks(doc.dom, ZONE_SELECT)
}

func (doc *PCDoc) GetAdvLinks() ([]string, error) {
	return GetLinks(doc.dom, ADV_SELECT)
}
