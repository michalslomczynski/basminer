package main

import (
	"fmt"
	"github.com/michalslomczynski/bas-opencv/browser"
	"github.com/michalslomczynski/bas-opencv/metamask"
	"log"
	"time"
)

func main() {
	metamask.DownloadMetamask()

	b := browser.LaunchBrowserWithExtension(metamask.FileName)
	page := b.MustPage("https://beta.blockapescissors.com/").MustWaitLoad()
	_, err := page.Element("canvas")
	if err != nil {
		log.Fatal(err)
	}
	pages, _ := b.Pages()
	for i := 0; i < 100; i++ {
		fmt.Println(len(pages))
		time.Sleep(1 * time.Second)
	}
	time.Sleep(100000 * time.Second)
	//el.MustScreenshot("sct.png")


	//window := gocv.NewWindow("test")
	//for i := 0; i < 200; i++ {
	//	img, err := screenshot.CaptureDisplay(0)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	mat, err := gocv.ImageToMatRGBA(img)
	//	window.IMShow(mat)
	//	window.WaitKey(500)
	//}
}
