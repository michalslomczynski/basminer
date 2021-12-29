package app

import (
	"context"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/michalslomczynski/bas-opencv/cvutil"
	"github.com/michalslomczynski/bas-opencv/util"
	"github.com/pkg/errors"
	"time"
)

const (
	roundLoopTimeout = 5 * time.Minute
	roundTimeout     = 1 * time.Second
	// Paths
	finishButtonPath     = prefix + "finish.png"
	exitGameButtonPath   = prefix + "exit_game.png"
	stoneButtonPath      = prefix + "stone.png"
	battleModeBannerPath = prefix + "battle_mode_banner.png"
	doneButtonPath       = prefix + "done_button.png"
	// Accuracies
	battleModeBannerAcc = 0.01
	stoneButtonAcc      = 0.03
	exitGameButtonAcc   = 0.03
	finishButtonAcc     = 0.03
	doneButtonAcc       = 0.03
	readyBoxAcc         = 0.07
)

func PlayRound(canvas *rod.Element) error {
	ctx, cancel := context.WithTimeout(context.Background(), roundLoopTimeout)
	defer cancel()

	for {
		findAndClickWithoutRetry(canvas, readyBoxPath, readyBoxAcc)
		findAndClickWithoutRetry(canvas, stoneButtonPath, stoneButtonAcc)
		findAndClickWithoutRetry(canvas, finishButtonPath, finishButtonAcc)
		findAndClickWithoutRetry(canvas, doneButtonPath, doneButtonAcc)
		exit := findAndClickWithoutRetry(canvas, exitGameButtonPath, exitGameButtonAcc)
		if exit {
			return nil
		}

		select {
		case <-ctx.Done():
			return errors.Wrapf(ctx.Err(), "round aborted")
		default:
			continue
		}
	}
}

func findFinishButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, finishButtonPath, mode, 0.01, lobbyTimeout)
}

func findAndClickWithoutRetry(canvas *rod.Element, path string, accuracy float32) bool {
	x, y, acc, err := cvutil.FindElementGenericWithoutRetry(canvas, path, mode)
	if err == nil {
		fmt.Printf("found %v with accuracy %v\n", path, acc)
		if mode < 2 {
			if acc < accuracy {
				for i := 0; i < 3; i++ {
					util.Click(canvas, x, y)
					time.Sleep(time.Millisecond * 300)
				}
				return true
			}
		} else {
			if acc > accuracy {
				for i := 0; i < 3; i++ {
					util.Click(canvas, x, y)
					time.Sleep(time.Millisecond * 300)
				}
				return true
			}
		}
	}
	return false
}

func findExitGameButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, exitGameButtonPath, mode, 0.01, lobbyTimeout)
}

func findStoneButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, stoneButtonPath, mode, 0.01, lobbyTimeout)
}

func findBattleModeBanner(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, battleModeBannerPath, mode, 0.01, lobbyTimeout)
}
