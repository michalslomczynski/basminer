package app

import (
	"github.com/go-rod/rod"
	"github.com/michalslomczynski/bas-opencv/cvutil"
	"github.com/michalslomczynski/bas-opencv/util"
)

const (
	battleModeButtonPath = prefix + "battle_mode.png"
	battleModeButtonAcc  = 0.03
)

func GoToLobby(canvas *rod.Element) error {
	x, y, err := findBattleModeButton(canvas)
	if err != nil {
		return err
	}

	util.Click(canvas, x, y)

	return nil
}

func findBattleModeButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, battleModeButtonPath, mode, 0.05, timeout)
}
