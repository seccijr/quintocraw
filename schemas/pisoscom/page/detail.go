package page

import (
	"github.com/seccijr/quintocrawl/model"
	"github.com/PuerkitoBio/goquery"
	"errors"
)

const THUMB_SELECT = ".frame.slideShow img"

func getThumbnails(dom *goquery.Document) []model.ImgNode {
	var result []model.ImgNode

	dom.Find(THUMB_SELECT).Each(func(i int, s *goquery.Selection) {
		imgNode := model.ToImgNode(s)
		result = append(result, imgNode)
	})

	return result
}

func photFromThumbnail(thumbnail *model.ImgNode) (model.ImgNode, error) {
	photo := &model.ImgNode{}
	if thumbnail.Src == "" {
		return nil, errors.New("No source for thumbnail")
	}

	return photo
}

func getPhotos(dom *goquery.Document) ([]model.ImgNode, error) {
	var result []model.ImgNode

	if Has(dom, THUMB_SELECT) {
		thumbnails := getThumbnails(dom)
		for thumbnail := range thumbnails {
			if photo, err := photFromThumbnail(&thumbnail); err == nil {
				result = append(result, photo)
			}
		}
	} else {
		return nil, errors.New("No photos nor thumbnails")
	}

	return result, nil
}

func (doc *PCDoc) ParseDetail() {

}