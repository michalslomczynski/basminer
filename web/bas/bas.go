package bas

import (
	"github.com/go-rod/rod"
	"github.com/michalslomczynski/bas-opencv/config"
	"github.com/pkg/errors"
)

const (
	appSelector = "canvas"
)

func OpenBASPage(b *rod.Browser) (*rod.Page, error) {
	var page *rod.Page
	err := rod.Try(func() {
		page = b.MustPage(config.Cfg.AppUrl).MustWaitLoad()
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open BAS page")
	}
	page.MustWaitLoad()
	return page, nil
}

func GetCanvas(p *rod.Page) (*rod.Element, error) {
	var appElem *rod.Element
	err := rod.Try(func() {
		appElem = p.MustElement(appSelector).MustWaitLoad()
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find appSelector element")
	}
	return appElem, nil
}
