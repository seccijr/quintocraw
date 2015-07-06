package page

import (
	"github.com/seccijr/quintocrawl/model"
	"github.com/PuerkitoBio/goquery"
	"errors"
	"regexp"
)

const THUMB_SELECT = ".frame.slideShow img"
const PHOTO_URL_PRE = "http://fotos.imghs.net/"
const PHOTO_REGEXP = "(s|m|l|xl)"

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
	photo := &model.ImgNode{}
	if thumb.Src == "" {
		return nil, errors.New("No source for thumb")
	}
	photo.Src = tSizePhotoUrl(thumb.Src, "l")

	return photo, nil
}

func getPhotos(dom *goquery.Document) ([]model.ImgNode, error) {
	var result []model.ImgNode

	if Has(dom, THUMB_SELECT) {
		thumbs := getThumbs(dom)
		result = getPhotosFromThumbs(thumbs)
	} else {
		return nil, errors.New("No photos nor thumbs")
	}

	return result, nil
}

func getPhotosFromThumbs(thumbs []model.ImgNode) ([]model.ImgNode, error) {
	var result []model.ImgNode

	for thumb := range thumbs {
		if photo, err := photoFromThumb(&thumb); err == nil {
			result = append(result, photo)
		}
	}

	return result, nil
}

func (doc *PCDoc) ParseDetail() {
	thumbs := getThumbs(doc.dom)
	photos, _ := getPhotosFromThumbs(thumbs)

}