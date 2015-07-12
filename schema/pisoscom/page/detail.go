package page

import (
	"github.com/seccijr/quintocrawl/model"
	"github.com/PuerkitoBio/goquery"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
	"strconv"
	"time"
	"fmt"
)

const ID_SELECT = "[name*='IdPiso']"
const THUMB_SELECT = ".frame.slideShow img"
const PHOTO_URL_PRE = "http://fotos.imghs.net/"
const PHOTO_REGEXP = "(s|m|l|xl)"
const TELF_ENC_SELECT = "[id='tlfEnc']"
const TELF_TXT_SELECT = ".number.one"
const INMO_SELECT = ".line.noMargin a[href^='/inmobiliaria']"
const DESC_BOD_SELECT = ".descriptionBlock .description"
const DATA_SELECT = "div.block"
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
const DATA_AREA_REGEXP = "((\\d+) m² construidos)?( / )?((\\d+) m² útiles)?"
const DATA_MORE_AGE_REGEXP = "más de (\\d+) años"
const DATA_LESS_AGE_REGEXP = "menos de (\\d+) años"

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

func basicDetails(result map[string]interface{}, dom *goquery.Selection) map[string]interface{} {
	re := regexp.MustCompile(DATA_LINE_REGEXP)
	dom.Find(DATA_LINE_SELECT).Each(func(i int, s *goquery.Selection) {
		match := re.FindStringSubmatch(s.Text())
		if match[1] != "" { result[match[1]] = match[2] }
	})

	return result
}

func listDetails(result map[string]interface{}, dom *goquery.Selection, key string) map[string]interface{} {
	details := make([]string, 10)
	dom.Find(DATA_LINE_SELECT).Each(func(i int, s *goquery.Selection) {
		details = append(details, s.Text())
	})
	if len(details) > 0 {
		result[key] = details
	}

	return result
}

func details(dom *goquery.Document) map[string]interface{} {
	result := make(map[string]interface{})
	dom.Find(DATA_SELECT).Each(func(num int, s *goquery.Selection) {
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

func convAge(age string) (time.Time, error) {
	var match []string
	var t time.Time
	reMore := regexp.MustCompile(DATA_MORE_AGE_REGEXP)
	reLess := regexp.MustCompile(DATA_LESS_AGE_REGEXP)
	switch {
	case reMore.MatchString(age):
		match = reMore.FindStringSubmatch(age)
		year, _ := strconv.Atoi(match[1])
		t = time.Now().AddDate(-year, 0, -1)
	case reLess.MatchString(age):
		match = reLess.FindStringSubmatch(age)
		year, _ := strconv.Atoi(match[1])
		t = time.Now().AddDate(-year, 0, 1)
	default:
		return t, errors.New("No matching time regexp")
	}

	return t, nil
}

func convArea(areaStr string) (model.Area, error) {
	var area model.Area
	re := regexp.MustCompile(DATA_AREA_REGEXP)
	match := re.FindStringSubmatch(areaStr)
	fmt.Println(match)

	return area, nil
}

func mapDetails(flat model.Flat, raw map[string]interface{}) model.Flat {
	re := regexp.MustCompile(DATA_N_REGEXP)
	if areaInt, exists := raw[DATA_AREA_KEY]; exists {
		if areaStr, ok := areaInt.(string); ok {
			area, _ := convArea(areaStr)
			flat.Area = area
		}
	}
	if roomsInt, exists := raw[DATA_ROOMS_KEY]; exists {
		if roomsStr, ok := roomsInt.(string); ok {
			match := re.FindString(roomsStr)
			n, _ := strconv.Atoi(match)
			flat.Rooms = n
		}
	}
	if bathsInt, exists := raw[DATA_BATHS_KEY]; exists {
		if bathsStr, ok := bathsInt.(string); ok {
			match := re.FindString(bathsStr)
			n, _ := strconv.Atoi(match)
			flat.Bathrooms = n
		}
	}
	if ageInt, exists := raw[DATA_AGE_KEY]; exists {
		if ageStr, ok := ageInt.(string); ok {
			n, _ := convAge(ageStr)
			flat.Age = n
		}
	}
	if maintInt, exists := raw[DATA_MAINT_KEY]; exists {
		if maintStr, ok := maintInt.(string); ok {
			flat.Maintenance = maintStr
		}
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
	details := details(doc.dom)
	fmt.Println(details)
	flat = mapDetails(flat, details)

	return flat, nil
}