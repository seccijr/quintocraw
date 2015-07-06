package page

import (
	"github.com/seccijr/quintocrawl/model"
	"github.com/PuerkitoBio/goquery"
	"encoding/base64"
	"errors"
	"regexp"
)

const THUMB_SELECT = ".frame.slideShow img"
const PHOTO_URL_PRE = "http://fotos.imghs.net/"
const PHOTO_REGEXP = "(s|m|l|xl)"
const TELF_SELECT = "#tlfEnc"
const INMO_SELECT = ".line.noMargin a[href^='/inmobiliaria']"
const DESC_BOD_SELECT = ".descriptionBlock .description"

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

func getPhotosFromThumbs(thumbs []model.ImgNode) ([]model.ImgNode, error) {
	var result []model.ImgNode

	for _, thumb := range thumbs {
		if photo, err := photoFromThumb(&thumb); err == nil {
			result = append(result, photo)
		}
	}

	return result, nil
}

func decTelf(dom *goquery.Document) (string, error) {
	val, exists := dom.Find(TELF_SELECT).First().Attr("value")
	if !exists {
		return "", errors.New("No matching telephone")
	}
	telByte, err := base64.StdEncoding.DecodeString(val)
	tel := string(telByte[:])
	if err == nil {
		return "", errors.New("Not Base64 telephone")
	}

	return tel, nil
}

func descBody(dom *goquery.Document) string {
	return dom.Find(DESC_BOD_SELECT).First().Text()
}

func (doc *PCDoc) ParseDetail() (model.Flat, error) {
	flat := model.Flat{}
	thumbs := getThumbs(doc.dom)
	photos, _ := getPhotosFromThumbs(thumbs)
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