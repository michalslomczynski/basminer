package app

import (
	"github.com/go-rod/rod"
	"github.com/michalslomczynski/bas-opencv/browser"
	"github.com/michalslomczynski/bas-opencv/metamask"
	"github.com/michalslomczynski/bas-opencv/web/bas"
	"github.com/michalslomczynski/bas-opencv/web/wallet"
	"log"
)

func Setup() (*rod.Element, error) {
	err := metamask.DownloadMetamask()
	if err != nil {
		return nil, err
	}

	b := browser.LaunchBrowserWithExtension(metamask.FileName, false)

	walletPage, err := wallet.ConnectToWallet(b)
	if err != nil {
		return nil, err
	}
	_ = walletPage

	BASPage, err := bas.OpenBASPage(b)
	if err != nil {
		return nil, err
	}

	canvas, err := bas.GetCanvas(BASPage)
	if err != nil {
		return nil, err
	}

	err = AcceptTerms(canvas)
	if err != nil {
		// Proceed with this error
		log.Println(err)
	}

	err = Login(canvas)
	if err != nil {
		return nil, err
	}

	return canvas, nil
}

func Restart(canvas *rod.Element) (*rod.Element, error) {
	log.Println("restarting...")

	BASPage := canvas.Page()
	BASPage.Reload()
	BASPage.MustWaitLoad()

	canvas, err := bas.GetCanvas(BASPage)
	if err != nil {
		return nil, err
	}

	err = Login(canvas)
	if err != nil {
		return nil, err
	}

	return canvas, nil
}
