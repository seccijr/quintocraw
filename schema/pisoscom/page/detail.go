package page

import (
	"github.com/seccijr/quintocrawl/model"
	"github.com/PuerkitoBio/goquery"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
	"strconv"
	"golang.org/x/tools/go/exact"
	"time"
)

const ID_SELECT = "[name*='IdPiso']"
const THUMB_SELECT = ".frame.slideShow img"
const PHOTO_URL_PRE = "http://fotos.imghs.net/"
const PHOTO_REGEXP = "(s|m|l|xl)"
const TELF_ENC_SELECT = "[id='tlfEnc']"
const TELF_TXT_SELECT = ".number.one"
const INMO_SELECT = ".line.noMargin a[href^='/inmobiliaria']"
const DESC_BOD_SELECT = ".descriptionBlock .description"
const DATA_SELECT = "div.characteristics > div.column > div.block"
const DATA_TITL_BASIC = "Datos básicos"
const DATA_TITL_EQUIP = "Equipamiento e instalaciones"
const DATA_TITL_CERTIFY = "Certificado energético"
const DATA_TITL_EXT = "Exteriores"
const DATA_TITL_FURN = "Muebles y acabados"
const DATA_TITL_SELECT = "h5"
const DATA_LINE_SELECT = "div.line"
const DATA_LINE_REGEXP = "^(.*):(.*)"
const DATA_AREA_KEY = "Superficie"
const DATA_ROOMS_KEY = "Superficie"
const DATA_BATHS_KEY = "Superficie"
const DATA_AGE_KEY = "Antigüedad"
const DATA_MAINT_KEY = "Conservación"
const DATA_N_REGEXP = "\\d+"

func tSizePhotoUrl(url string, size string) string {
	r := regexp.MustCompile(PHOTO_URL_PRE + PHOTO_REGEXP)
	return r.ReplaceAllString(url, PHOTO_URL_PRE + size)
}

func getThumbs(dom *goquery.Document) []model.ImgNode {
	var result []model.ImgNode

	dom.Find(THUMB_SELECT).Each(func(i int, s *goquery.Selection) {
		imgNode := model.ToImgNode(s)
		result = append(result, imgNode)
	})

	return result
}

func photoFromThumb(thumb *model.ImgNode) (model.ImgNode, error) {
	photo := model.ImgNode{}
	if thumb.Src == "" {
		return photo, errors.New("No source for thumb")
	}
	photo.Src = tSizePhotoUrl(thumb.Src, "l")

	return photo, nil
}

func photosFromThumbs(thumbs []model.ImgNode) ([]model.ImgNode, error) {
	var result []model.ImgNode

	for _, thumb := range thumbs {
		if photo, err := photoFromThumb(&thumb); err == nil {
			result = append(result, photo)
		}
	}

	return result, nil
}

func decTelf(dom *goquery.Document) (string, error) {
	var tel string
	if Has(dom, TELF_ENC_SELECT) {
		val, exists := dom.Find(TELF_ENC_SELECT).First().Attr("value")
		if !exists {
			return "", errors.New("No matching encoded telephone")
		}
		telByte, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return "", err
		}
		tel = string(telByte[:])
	} else {
		tel = dom.Find(TELF_TXT_SELECT).First().Text()
	}

	return tel, nil
}

func descBody(dom *goquery.Document) string {
	return dom.Find(DESC_BOD_SELECT).First().Text()
}

func getRef(dom *goquery.Document) (string, error) {
	val, exists := dom.Find(ID_SELECT).First().Attr("value")
	if !exists {
		return "", errors.New("No matching ID")
	}

	return val, nil
}

func basicDetails(result map[string]interface{}, dom *goquery.Document) map[string]interface{} {
	re := regexp.MustCompile(DATA_LINE_REGEXP)
	dom.Find(DATA_LINE_SELECT).Each(func (i int, s *goquery.Selection) {
		match := re.FindStringSubmatch(s.Text())
		if match[1] != "" { result = append(result, match[2]) }
	})

	return result
}

func listDetails(result map[string]interface{}, key string, dom *goquery.Document) map[string]interface{} {
	result[key] = make([]string)
	dom.Find(DATA_LINE_SELECT).Each(func (i int, s *goquery.Selection) {
		result[key] = append(result[key], s.Text())
	})

	return result
}

func details(dom *goquery.Document) map[string]interface{} {
	result := make(map[string]interface{})
	dom.Find(DATA_SELECT).Each(func (num int, s *goquery.Selection) {
		title := s.Find(DATA_TITL_SELECT).First().Text()
		switch {
		case strings.Contains(title, DATA_TITL_BASIC):
			result = basicDetails(result, s)
		case strings.Contains(title, DATA_TITL_EQUIP):
			result = listDetails(result, s, DATA_TITL_EQUIP)
		case strings.Contains(title, DATA_TITL_EXT):
			result = listDetails(result, s, DATA_TITL_EXT)
		case strings.Contains(title, DATA_TITL_FURN):
			result = listDetails(result, s, DATA_TITL_FURN)
		}
	})

	return result
}

func convAge(age string) time.Time{

}

func mapDetails(flat model.Flat, raw map[string]interface{}) model.Flat {
	re := regexp.MustCompile(DATA_N_REGEXP)
	if area, exists := raw[DATA_AREA_KEY]; exists {
		flat.Area = area
	}
	if val, exists := raw[DATA_ROOMS_KEY]; exists {
		match := re.FindString(val)
		n, _ := strconv.Atoi(match)
		flat.Rooms = n
	}
	if val, exists := raw[DATA_BATHS_KEY]; exists {
		match := re.FindString(val)
		n, _ := strconv.Atoi(match)
		flat.Bathrooms = n
	}
	if val, exists := raw[DATA_AGE_KEY]; exists {
		match := re.FindString(val)
		n, _ := strconv.Atoi(match)
		flat.Bathrooms = n
	}
	if val, exists := raw[DATA_MAINT_KEY]; exists {
		match := re.FindString(val)
		n, _ := strconv.Atoi(match)
		flat.Bathrooms = n
	}

	return flat
}

func (doc *PCDoc) ParseDetail() (model.Flat, error) {
	flat := model.Flat{}
	ref, err := getRef(doc.dom)
	if err != nil {
		return flat, err
	}
	flat.Ref = ref
	thumbs := getThumbs(doc.dom)
	photos, _ := photosFromThumbs(thumbs)
	flat.Pictures = photos
	tel, err := decTelf(doc.dom)
	if err != nil {
		return flat, err
	}
	flat.Telephone = tel
	desc := descBody(doc.dom)
	flat.Description = desc

	return flat, nil
}