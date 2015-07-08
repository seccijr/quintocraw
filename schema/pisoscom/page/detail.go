package page

import (
	"github.com/seccijr/quintocrawl/model"
	"github.com/PuerkitoBio/goquery"
	"encoding/base64"
	"errors"
	"regexp"
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