package model
import "github.com/PuerkitoBio/goquery"

type ImgNode struct {
	Src     string
	Alt     string
	Comment string
}

func ToImgNode(s *goquery.Selection) ImgNode {
	imgNode := ImgNode{}
	if lazy, hasLazy := s.Attr("data-pslider-lazy"); hasLazy {
		imgNode.Src = lazy
	} else if src, hasSrc := s.Attr("src"); hasSrc {
		imgNode.Src = src
	}

	if alt, hasAlt := s.Attr("alt"); hasAlt {
		imgNode.Alt = alt
	}

	if comm, hasComm := s.Attr("data-comentario"); hasComm {
		imgNode.Comment = comm
	}

	return imgNode
}
