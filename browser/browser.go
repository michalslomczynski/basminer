package browser

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"path/filepath"
)

func LaunchBrowser() *rod.Browser {
	browser := rod.New().MustConnect()
	return browser
}

func LaunchHeadedBrowser() *rod.Browser {
	u := launcher.New().
		Headless(false).
		MustLaunch()
	return rod.New().ControlURL(u).MustConnect()
}

func LaunchBrowserWithExtension(path string) *rod.Browser {
	extPath, _ := filepath.Abs(path)
	fmt.Println(extPath)

	u := launcher.New().
		Set("load-extension", extPath).
		Headless(false).
		MustLaunch()
	return rod.New().ControlURL(u).MustConnect()
	//page.MustWait(`document.title === 'test-extension'`)
}
