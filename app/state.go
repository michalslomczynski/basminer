package app

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/michalslomczynski/bas-opencv/cvutil"
	cv "gocv.io/x/gocv"
)

type State struct {
	menu  bool
	lobby bool
	round bool
}

type StateAssets struct {
	menu   *cv.Mat
	lobby  *cv.Mat
	round  *cv.Mat
	offset *proto.DOMRect
}

// StateUpdate provides performant feedback about current state of image
func StateUpdate(image *cv.Mat, assets *StateAssets, state *State) {
	fmt.Printf("\nfound lobby ")
	state.lobby = cvutil.IsMatElementVisible(image, assets.lobby, assets.offset, mode, roomNameAcc)
	fmt.Printf("\nfound menu ")
	state.menu = cvutil.IsMatElementVisible(image, assets.menu, assets.offset, mode, battleModeButtonAcc)
	fmt.Printf("\nfound round ")
	state.round = cvutil.IsMatElementVisible(image, assets.round, assets.offset, mode, battleModeBannerAcc)
	fmt.Println()
}

func loadStateAssets(canvas *rod.Element) (*StateAssets, error) {
	assets := &StateAssets{}

	shape, err := canvas.Shape()
	if err != nil {
		return nil, err
	}
	assets.offset = shape.Box()

	assets.menu, err = cvutil.LoadAssetToMat(battleModeButtonPath)
	if err != nil {
		return nil, err
	}

	assets.lobby, err = cvutil.LoadAssetToMat(roomNamePath)
	if err != nil {
		return nil, err
	}

	assets.round, err = cvutil.LoadAssetToMat(battleModeBannerPath)
	if err != nil {
		return nil, err
	}

	return assets, nil
}
