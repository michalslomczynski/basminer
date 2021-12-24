package web

import "github.com/go-rod/rod"

const extensionUrl = "chrome-extension://nojmdacjahoombflmhcebljpcpjfogei/home.html#initialize/welcome"

func ConnectToWallet(browser *rod.Browser) {
	browser.MustPage(extensionUrl).MustWaitLoad()
	pressGetStarted()
}

func pressGetStarted(page ) {

}