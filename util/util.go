package util

import (
	"bytes"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/pkg/errors"
	"image"
	"path/filepath"
)

func BytesToImageSave(bytes []byte, path string) error {
	err := utils.OutputFile(filepath.Join(path), bytes)
	if err != nil {
		return errors.Wrapf(err, "failed to save bytes to image")
	}
	return nil
}

func ElemSctToBytes(elem *rod.Element) ([]byte, error) {
	if elem == nil {
		return nil, errors.New("element cannot be nil")
	}

	bytes, err := elem.Screenshot(proto.PageCaptureScreenshotFormatPng, 1)
	if err != nil {
		return nil, errors.Wrapf(err, "could not take screenshot")
	}

	return bytes, nil
}

func ElemSctToImage(elem *rod.Element) (*image.Image, error) {
	b, err := ElemSctToBytes(elem)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func Click(elem *rod.Element, x, y int) {
	elem.Page().Mouse.Move(float64(x), float64(y), 1)
	elem.Page().Mouse.Click(proto.InputMouseButtonLeft)
}

func ScrollDown(elem *rod.Element) error {
	return elem.Page().Mouse.Scroll(0, 100000, 1)
}
