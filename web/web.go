package web

import (
	"github.com/go-rod/rod"
	"time"
)

const url = "https://chrome.google.com/webstore/detail/metamask/nkbihfbeogaeaoehlefnkodbefgpgknn?hl=en"

func DownloadMetaMask(b *rod.Browser) {
	page := b.MustPage(url).MustWaitLoad()
	page.MustElement(".VfPpkd-RLmnJb").MustClick()

	time.Sleep(10000 * time.Second)
}
