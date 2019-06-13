package imagematch

import (
	"errors"
	"fmt"
	"image"
	"io"

	// must
	_ "image/jpeg"
	_ "image/png"

	"gocv.io/x/gocv"
)

const (
	e = 1E9
)

// Template of match
type Template struct {
	template gocv.Mat
	sill     float32
}

// NewTemplateFromStream only create template from io.Reader
func NewTemplateFromStream(r io.Reader, sill ...float32) (*Template, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	m, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return nil, err
	}
	var s float32
	if len(sill) > 0 {
		s = sill[0]
	} else {
		s = 0.95
	}
	return &Template{m, s}, nil
}

// Match picture
func (t Template) Match(img gocv.Mat) (bool, error) {
	var err error
	defer func() {
		if rerr := recover(); rerr != nil {
			err = fmt.Errorf("%v", rerr)
		}
	}()
	result := gocv.NewMat()
	defer result.Close()
	m := gocv.NewMat()
	gocv.MatchTemplate(img, t.template, &result, gocv.TmCcoeffNormed, m)
	m.Close()
	_, max, _, _ := gocv.MinMaxLoc(result)
	max = max / e
	return max < t.sill, err
}

// Close mat
func (t *Template) Close() {
	(&t.template).Close()
}

// NewMatFromFile if error, return emtpy gocv.Mat
func NewMatFromFile(path string) (gocv.Mat, error) {
	img := gocv.IMRead(path, gocv.IMReadAnyDepth)
	if img.Empty() {
		return img, errors.New("open picture error")
	}
	return img, nil
}

// NewMatFromStream reuse ReadSeeker
func NewMatFromStream(r io.Reader) (gocv.Mat, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return gocv.Mat{}, err
	}
	return gocv.ImageToMatRGB(img)
}

// ImageMatch read from file
func ImageMatch(template, target string, sill float32) (bool, error) {
	imgTemplate := gocv.IMRead(template, gocv.IMReadGrayScale)
	defer imgTemplate.Close()
	if imgTemplate.Empty() {
		return false, errors.New("open template picture error")
	}
	imgSence := gocv.IMRead(target, gocv.IMReadGrayScale)
	defer imgSence.Close()
	if imgSence.Empty() {
		return false, errors.New("open template picture error")
	}
	result := gocv.NewMat()
	defer result.Close()
	m := gocv.NewMat()
	gocv.MatchTemplate(imgSence, imgTemplate, &result, gocv.TmSqdiff, m)
	m.Close()
	_, max, _, _ := gocv.MinMaxLoc(result)
	max = max / (e)
	return max < sill, nil
}
