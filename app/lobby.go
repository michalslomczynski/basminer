package app

import (
	"context"
	"github.com/go-rod/rod"
	"github.com/go-vgo/robotgo"
	"github.com/michalslomczynski/bas-opencv/config"
	"github.com/michalslomczynski/bas-opencv/cvutil"
	"github.com/michalslomczynski/bas-opencv/util"
	"github.com/pkg/errors"
	"time"
)

const (
	lobbyLoopTimeout = 1 * time.Minute
	lobbyTimeout     = 10 * time.Second
	// Paths
	lobbyScrollBarPath     = prefix + "scroll_bar.png"
	joinButtonPath         = prefix + "join_button.png"
	joinFailedOkButtonPath = prefix + "join_failed_ok.png"
	newGameButtonPath      = prefix + "new_game.png"
	gameNameInputPath      = prefix + "game_name.png"
	startGameButtonPath    = prefix + "start_game.png"
	readyBoxPath           = prefix + "ready_box.png"
	loadingPath            = prefix + "loading.png"
	roomNamePath           = prefix + "room_name.png"
	// Accuracies
	lobbyJoinButtonAcc     = 0.04
	lobbyFailedOkButtonAcc = 0.04
	roomNameAcc            = 0.06
	readyBoxAcc            = 0.02
)

func SelectRoom(canvas *rod.Element) error {
	ctx, cancel := context.WithTimeout(context.Background(), lobbyLoopTimeout)
	defer cancel()

	m := canvas.Page().Mouse
	for {
		x, y, err := findJoinButton(canvas)
		if err == nil {
			// Because client is laggy and unresponsive
			for j := 0; j < 4; j++ {
				util.Click(canvas, x, y)
				time.Sleep(time.Millisecond * 300)
			}
			_, _, err = findLoading(canvas)
			if err == nil {
				break
			}
		}

		// Try to scroll down list of rooms because of known server bug
		for i := 0; i < 15; i++ {
			m.Move(float64(x), float64(y), 1)
			m.Scroll(0, 10000000, 1)
		}

		x, y, err = findJoinFailedOkButton(canvas)
		if err == nil {
			// Because client is laggy and unresponsive
			for j := 0; j < 4; j++ {
				util.Click(canvas, x, y)
				time.Sleep(time.Millisecond * 300)
			}
		}

		select {
		case <-ctx.Done():
			return errors.Wrapf(ctx.Err(), "room select aborted")
		default:
			continue
		}
	}

	return nil
}

func NewGame(canvas *rod.Element) error {
	x, y, err := findNewGameButton(canvas)
	if err != nil {
		return err
	}
	util.Click(canvas, x, y)

	x, y, err = findGameNameInput(canvas)
	if err != nil {
		return err
	}
	util.Click(canvas, x, y)

	canvas.Page().Keyboard.MustInsertText(config.Cfg.LobbyName)
	robotgo.TypeStr(config.Cfg.LobbyName)

	x, y, err = findStartGameButton(canvas)
	if err != nil {
		return err
	}
	util.Click(canvas, x, y)

	x, y, err = findReadyBox(canvas)
	if err != nil {
		return err
	}
	util.Click(canvas, x, y)

	return nil
}

func findLobbyScrollButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, lobbyScrollBarPath, mode, 0.05, timeout)
}

func findJoinButton(canvas *rod.Element) (int, int, error) {
	canvas.MustScreenshot("debugSct.png")
	return cvutil.FindElementGeneric(canvas, joinButtonPath, mode, lobbyJoinButtonAcc, lobbyTimeout)
}

func findJoinFailedOkButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, joinFailedOkButtonPath, mode, lobbyFailedOkButtonAcc, lobbyTimeout)
}

func findNewGameButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, newGameButtonPath, mode, 0.01, timeout)
}

func findGameNameInput(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, gameNameInputPath, mode, 0.01, timeout)
}

func findStartGameButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, startGameButtonPath, mode, 0.01, timeout)
}

func findReadyBox(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, readyBoxPath, mode, 0.01, lobbyTimeout)
}

func findAndClickReadyBoxWithoutRetry(canvas *rod.Element) {
	x, y, acc, err := cvutil.FindElementGenericWithoutRetry(canvas, readyBoxPath, mode)
	if err == nil {
		if mode < 2 {
			if acc < readyBoxAcc {
				util.Click(canvas, x, y)
			}
		} else {
			if acc > readyBoxAcc {
				util.Click(canvas, x, y)
			}
		}
	}
}

func findLoading(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, loadingPath, mode, 0.01, lobbyTimeout)
}

func findRoomName(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, roomNamePath, mode, 0.01, lobbyTimeout)
}
