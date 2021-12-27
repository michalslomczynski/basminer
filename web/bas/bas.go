package bas

import (
	"github.com/go-rod/rod"
	"github.com/pkg/errors"
)

const (
	appUrl = "https://beta.blockapescissors.com"
	// Selectors
	app = "canvas"
)

func OpenBASPage(b *rod.Browser) (*rod.Page, error) {
	var page *rod.Page
	err := rod.Try(func() {
		page = b.MustPage(appUrl).MustWaitLoad()
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
		appElem = p.MustElement(app).MustWaitLoad()
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find app element")
	}
	return appElem, nil
}
