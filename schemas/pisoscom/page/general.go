package page

const STATE_SELECT = ".zonas.clearfix a.bold"
const PROV_SELECT = ".zonas.clearfix a[class='']"
const ZONE_SELECT = ".zoneList a"
const ADV_SELECT = ".gridList a.anuncioLink"
const DETAIL_SELECT = ".content.Detalle"

func (doc *PCDoc) HasStates() bool {
	return Has(doc.dom, STATE_SELECT)
}

func (doc *PCDoc) HasProvinces() bool {
	return Has(doc.dom, PROV_SELECT)
}

func (doc *PCDoc) HasZones() bool {
	return Has(doc.dom, ZONE_SELECT)
}

func (doc *PCDoc) HasAds() bool {
	return Has(doc.dom, ADV_SELECT)
}

func (doc *PCDoc) IsDetail() bool {
	return Has(doc.dom, DETAIL_SELECT)
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
