package page
import (
	"regexp"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"fmt"
	"errors"
)

const STATE_SELECT = ".zonas.clearfix a.bold"
const PROV_SELECT = ".zonas.clearfix a[class='']"
const ZONE_SELECT = ".zoneList a"
const ADV_SELECT = ".gridList a.anuncioLink"
const DETAIL_SELECT = ".content.Detalle"
const PAG_NUM_REGEXP = "([0-9]+)/$"
const PAG_LAST_SELECT = ".pager a.item:nth-last-child(2)"

func pageNum(url string) string {
	re := regexp.MustCompile(PAG_NUM_REGEXP)
	m := re.FindStringSubmatch(url)

	return m[1]
}

func maxPages(dom *goquery.Document) (int, error) {
	last, err := lastPageLink(dom)
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(pageNum(last))
}

func pageLinks(dom *goquery.Document) ([]string, error) {
	var links []string
	n, err := maxPages(dom)
	if err != nil {
		return links, err
	}
	url, err := lastPageLink(dom)
	if err != nil {
		return links, err
	}
	plain := removePageNumber(url)
	for i:= 1; i <= n; i++ {
		links = append(links, fmt.Sprintf("%s/%d/", plain, n))
	}

	return links, nil
}

func removePageNumber(url string) string {
	re := regexp.MustCompile(PAG_NUM_REGEXP)
	return re.ReplaceAllString(url, "")
}

func lastPageLink(dom *goquery.Document) (string, error) {
	if !Has(dom, PAG_LAST_SELECT) {
		return "", errors.New("Without last page")
	}

	last, exists := dom.Find(PAG_LAST_SELECT).First().Attr("href")
	if !exists {
		return "", errors.New("Without href in last page")
	}

	return last, nil
}

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
	var result []string

	pages, _ := pageLinks(doc.dom)
	ads, err := GetLinks(doc.dom, ADV_SELECT)
	if err != nil {
		return result, err
	}
	result = append(ads, pages...)

	return result, nil
}
