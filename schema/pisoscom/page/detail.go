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
	"math"
)

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
	var t time.Time
	switch {
	case strings.Contains(age, DATA_LESS_5_Y_REGEXP):
		t = time.Now()
	case strings.Contains(age, DATA_B_5_10_Y_REGEXP):
		t = time.Now().AddDate(-5, 0, -1)
	case strings.Contains(age, DATA_B_10_20_Y_REGEXP):
		t = time.Now().AddDate(-10, 0, -1)
	case strings.Contains(age, DATA_B_20_30_Y_REGEXP):
		t = time.Now().AddDate(-20, 0, -1)
	case strings.Contains(age, DATA_B_30_50_Y_REGEXP):
		t = time.Now().AddDate(-30, 0, -1)
	case strings.Contains(age, DATA_MORE_50_Y_REGEXP):
		t = time.Now().AddDate(-50, 0, -1)
	default:
		return t, errors.New("No matching time regexp")
	}

	return t, nil
}

func convFees(fee string) (model.PriceRange, error) {
	var r model.PriceRange
	from := model.Price{Currency: "€"}
	to := model.Price{Currency: "€"}
	switch {
	case strings.Contains(fee, DATA_B_10_20_E_REGEXP):
		from.Amount = 10.0
		to.Amount = 20.0
	case strings.Contains(fee, DATA_B_20_40_E_REGEXP):
		from.Amount = 20.0
		to.Amount = 40.0
	case strings.Contains(fee, DATA_B_40_60_E_REGEXP):
		from.Amount = 40.0
		to.Amount = 60.0
	case strings.Contains(fee, DATA_B_60_80_E_REGEXP):
		from.Amount = 60.0
		to.Amount = 80.0
	case strings.Contains(fee, DATA_B_80_100_E_REGEXP):
		from.Amount = 80.0
		to.Amount = 100.0
	case strings.Contains(fee, DATA_MORE_100_E_REGEXP):
		from.Amount = 100.0
		to.Amount = math.Inf(1)
	default:
		return r, errors.New("No matching time regexp")
	}

	r = model.PriceRange{from, to}

	return r, nil
}

func convArea(areaStr string) (model.Area, error) {
	var area model.Area
	re := regexp.MustCompile(DATA_AREA_REGEXP)
	match := re.FindStringSubmatch(areaStr)
	built, err := strconv.ParseFloat(match[1], 64)
	if err == nil {
		area.Built = built
	}
	util, err := strconv.ParseFloat(match[4], 64)
	if err == nil {
		area.Util = util
	}

	return area, nil
}

func mapIntDetail(raw map[string]interface{}, key string) int {
	var result int
	re := regexp.MustCompile(DATA_N_REGEXP)
	if roomsInt, exists := raw[key]; exists {
		if roomsStr, ok := roomsInt.(string); ok {
			match := re.FindString(roomsStr)
			n, _ := strconv.Atoi(match)
			result = n
		}
	}

	return result
}

func mapListDetail(raw map[string]interface{}, key string) []string {
	var result []string

	if inter, exists := raw[key]; exists {
		if list, ok := inter.([]string); ok {
			result = list
		}
	}

	return result
}

func mapDetails(flat model.Flat, raw map[string]interface{}) model.Flat {
	if areaInt, exists := raw[DATA_AREA_KEY]; exists {
		if areaStr, ok := areaInt.(string); ok {
			area, _ := convArea(areaStr)
			flat.Area = area
		}
	}
	flat.Rooms = mapIntDetail(raw, DATA_ROOMS_KEY)
	flat.Bathrooms = mapIntDetail(raw, DATA_BATHS_KEY)
	flat.Floor  = mapIntDetail(raw, DATA_FLOOR_KEY)
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
	if feeInt, exists := raw[DATA_COM_FEES_KEY]; exists {
		if feeStr, ok := feeInt.(string); ok {
			n, _ := convFees(feeStr)
			flat.ComFees = n
		}
	}
	flat.Equipment = mapListDetail(raw, DATA_TITL_EQUIP)
	flat.Exterior = mapListDetail(raw, DATA_TITL_EXT)
	flat.Furniture = mapListDetail(raw, DATA_TITL_FURN)

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
	flat = mapDetails(flat, details)
	fmt.Println(flat)
	return flat, nil
}