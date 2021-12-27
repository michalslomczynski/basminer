package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"log"
	"path/filepath"
)

const (
	width  = 1280
	height = 800
)

func LaunchBrowser() *rod.Browser {
	browser := rod.New().MustConnect()
	return browser
}

func LaunchBrowserHeaded() *rod.Browser {
	u := launcher.New().
		Headless(false).
		MustLaunch()
	return rod.New().ControlURL(u).MustConnect()
}

func LaunchBrowserWithExtension(path string, headless bool) *rod.Browser {
	extPath, _ := filepath.Abs(path)

	u := launcher.New().
		Set("load-extension", extPath).
		Headless(headless).
		MustLaunch()
	return rod.New().ControlURL(u).MustConnect()
}

func GetTargetByTitle(b *rod.Browser, title string) *proto.TargetTargetInfo {
	list, _ := proto.TargetGetTargets{}.Call(b)
	for _, target := range list.TargetInfos {
		if target.Title == title {
			return target
		}
	}
	return nil
}

func GetTargetByURL(b *rod.Browser, url string) *proto.TargetTargetInfo {
	list, _ := proto.TargetGetTargets{}.Call(b)
	for _, target := range list.TargetInfos {
		if target.URL == url {
			return target
		}
	}
	return nil
}

func GetPageFromTargetTitle(b *rod.Browser, title string) *rod.Page {
	target := GetTargetByTitle(b, title)
	if target == nil {
		return nil
	}
	if target.Type != proto.TargetTargetInfoTypePage {
		log.Printf("could not retrieve page from target of type: %v", target.Type)
		return nil
	}
	page, err := b.PageFromTarget(target.TargetID)
	if err != nil {
		return nil
	}
	return page
}

// InnerAbsolutePosition returns absolute document position on the screen.
func InnerAbsolutePosition(page *rod.Page) (int, int, error) {
	res, err := page.Eval("window.outerHeight - window.innerHeight + window.screenTop")
	if err != nil {
		return 0, 0, err
	}
	y := int(res.Value.Val().(float64))

	res, err = page.Eval("window.screenLeft")
	if err != nil {
		return 0, 0, err
	}
	x := int(res.Value.Val().(float64))

	return x, y, nil
}
