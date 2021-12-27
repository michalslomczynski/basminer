package cvutil

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/michalslomczynski/bas-opencv/util"
	"github.com/pkg/errors"
	cv "gocv.io/x/gocv"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

func GetMatchLocation(img, tmpl *cv.Mat, mode cv.TemplateMatchMode) (int, int, float32) {
	res := cv.NewMat()
	cv.MatchTemplate(*img, *tmpl, &res, mode, cv.NewMat())

	minVal, maxVal, minLoc, maxLoc := cv.MinMaxLoc(res)

	// See opencv docs matchTemplate() for reference
	if mode < 2 {
		return minLoc.X, minLoc.Y, minVal
	} else {
		return maxLoc.X, maxLoc.Y, maxVal
	}
}

func LoadAssetToMat(path string) (*cv.Mat, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	mat, err := BytesToMat(b)
	if err != nil {
		return nil, err
	}

	return mat, nil
}

func BytesToMat(b []byte) (*cv.Mat, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode image from bytes")
	}
	mat, err := cv.ImageToMatRGBA(img)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert image to Mat")
	}
	return &mat, nil
}

func SaveImgFromBytes(b []byte, name string) error {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}

	out, err := os.Create(name)
	defer out.Close()
	if err != nil {
		return err
	}

	err = png.Encode(out, img)
	if err != nil {
		return err
	}
	return nil
}

func FindElementGeneric(canvas *rod.Element, templatePath string, mode cv.TemplateMatchMode, accuracy float32, timeout time.Duration) (int, int, error) {
	shape, err := canvas.Shape()
	if err != nil {
		return 0, 0, err
	}

	// Contains relative position of canvas on page
	rb := shape.Box()

	tmpl, err := LoadAssetToMat(templatePath)
	if err != nil {
		return 0, 0, err
	}

	// Template size offset values
	tsz := tmpl.Size()
	tx := tsz[0] / 2
	ty := tsz[0] / 2

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		bb, err := util.ElemSctToBytes(canvas)
		if err != nil {
			return 0, 0, err
		}

		img, err := BytesToMat(bb)
		if err != nil {
			return 0, 0, err
		}

		x, y, acc := GetMatchLocation(img, tmpl, mode)
		fmt.Printf("found template %s at x=%v y=%v with accuracy %v\n", templatePath, x, y, acc)

		if mode < 2 {
			if acc < accuracy {
				return x + int(rb.X) + tx, y + int(rb.Y) + ty, nil
			}
		} else {
			if acc > accuracy {
				return x + int(rb.X) + tx, y + int(rb.Y) + ty, nil
			}
		}

		select {
		case <-ctx.Done():
			return 0, 0, errors.Wrapf(ctx.Err(), "could not find element from %v", templatePath)
		default:
			time.Sleep(time.Second)
			continue
		}
	}
}

func FindElementGenericWithoutRetry(canvas *rod.Element, templatePath string, mode cv.TemplateMatchMode) (int, int, float32, error) {
	shape, err := canvas.Shape()
	if err != nil {
		return 0, 0, 0, err
	}

	// Contains relative position of element on page
	rb := shape.Box()

	tmpl, err := LoadAssetToMat(templatePath)
	if err != nil {
		return 0, 0, 0, err
	}

	// Template size offset values
	tsz := tmpl.Size()
	tx := tsz[0] / 2
	ty := tsz[0] / 2

	img, err := ElemToMat(canvas)
	if err != nil {
		return 0, 0, 0, err
	}

	x, y, acc := GetMatchLocation(img, tmpl, mode)

	return x + int(rb.X) + tx, y + int(rb.Y) + ty, acc, nil
}

func FindMatElement(img, tmpl *cv.Mat, offset *proto.DOMRect, mode cv.TemplateMatchMode) (int, int, float32) {
	tsz := tmpl.Size()
	tx := tsz[0] / 2
	ty := tsz[0] / 2

	x, y, acc := GetMatchLocation(img, tmpl, mode)

	return x + int(offset.X) + tx, y + int(offset.Y) + ty, acc
}

func IsMatElementVisible(img, tmpl *cv.Mat, offset *proto.DOMRect, mode cv.TemplateMatchMode, accuracy float32) bool {
	_, _, acc := FindMatElement(img, tmpl, offset, mode)
	fmt.Printf("found mat element with accuracy %v\n", acc)
	if mode < 2 {
		if acc < accuracy {
			return true
		}
	} else {
		if acc > accuracy {
			return true
		}
	}

	return false
}

func ElemToMat(elem *rod.Element) (*cv.Mat, error) {
	bb, err := util.ElemSctToBytes(elem)
	if err != nil {
		return nil, err
	}

	img, err := BytesToMat(bb)
	if err != nil {
		return nil, err
	}

	return img, nil
}
