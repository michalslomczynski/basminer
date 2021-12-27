package app

import (
	"context"
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
	// Accuracies
	battleModeBannerAcc = 0.01
	stoneButtonAcc      = 0.03
	exitGameButtonAcc   = 0.03
)

func PlayRound(canvas *rod.Element) error {
	ctx, cancel := context.WithTimeout(context.Background(), roundLoopTimeout)
	defer cancel()

	for {
		findAndClickReadyBoxWithoutRetry(canvas)

		findAndClickStoneButtonWithoutRetry(canvas)

		findAndClickFinishButtonWithoutRetry(canvas)

		exit := findAndClickExitGameButtonWithoutRetry(canvas)
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

func findAndClickFinishButtonWithoutRetry(canvas *rod.Element) bool {
	x, y, acc, err := cvutil.FindElementGenericWithoutRetry(canvas, finishButtonPath, mode)
	if err == nil {
		if mode < 2 {
			if acc < exitGameButtonAcc {
				for i := 0; i < 3; i++ {
					util.Click(canvas, x, y)
					time.Sleep(time.Millisecond * 300)
				}
				return true
			}
		} else {
			if acc > exitGameButtonAcc {
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

func findAndClickExitGameButtonWithoutRetry(canvas *rod.Element) bool {
	x, y, acc, err := cvutil.FindElementGenericWithoutRetry(canvas, exitGameButtonPath, mode)
	if err == nil {
		if mode < 2 {
			if acc < exitGameButtonAcc {
				for i := 0; i < 3; i++ {
					util.Click(canvas, x, y)
					time.Sleep(time.Millisecond * 300)
				}
				return true
			}
		} else {
			if acc > exitGameButtonAcc {
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

func findStoneButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, stoneButtonPath, mode, 0.01, lobbyTimeout)
}

func findAndClickStoneButtonWithoutRetry(canvas *rod.Element) bool {
	x, y, acc, err := cvutil.FindElementGenericWithoutRetry(canvas, stoneButtonPath, mode)
	if err == nil {
		if mode < 2 {
			if acc < stoneButtonAcc {
				for i := 0; i < 3; i++ {
					util.Click(canvas, x, y)
					time.Sleep(time.Millisecond * 300)
				}
				return true
			}
		} else {
			if acc > stoneButtonAcc {
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

func findBattleModeBanner(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, battleModeBannerPath, mode, 0.01, lobbyTimeout)
}
