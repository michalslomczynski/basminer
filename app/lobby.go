package app

import (
	"context"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-vgo/robotgo"
	"github.com/michalslomczynski/bas-opencv/config"
	"github.com/michalslomczynski/bas-opencv/cvutil"
	"github.com/michalslomczynski/bas-opencv/util"
	"github.com/pkg/errors"
	"time"
)

const (
	lobbyLoopTimeout = 1 * time.Minute
	lobbyTimeout     = 4 * time.Second
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
	lobbyScrollSpecialPath = prefix + "lobby_scroll_special.png"
	// Accuracies
	lobbyScrollAcc         = 0.08
	lobbyScrollSpecialAcc  = 0.09
	lobbyJoinButtonAcc     = 0.06
	lobbyFailedOkButtonAcc = 0.05
	roomNameAcc            = 0.06
)

func SelectRoomAlt(canvas *rod.Element) error {
	findAndClickWithoutRetry(canvas, joinButtonPath, lobbyJoinButtonAcc)
	findAndClickWithoutRetry(canvas, joinFailedOkButtonPath, lobbyFailedOkButtonAcc)

	return nil
}

func SelectRoom(canvas *rod.Element) error {
	ctx, cancel := context.WithTimeout(context.Background(), lobbyLoopTimeout)
	defer cancel()

	for {
		x, y, err := findJoinButton(canvas)
		if err == nil {
			// Because client is laggy and unresponsive
			for j := 0; j < 4; j++ {
				util.Click(canvas, x, y)
				time.Sleep(time.Millisecond * 300)
			}
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
	return cvutil.FindElementGeneric(canvas, lobbyScrollBarPath, mode, 0.05, lobbyTimeout)
}

func findAndClickLobbyScrollButtonWithoutRetry(canvas *rod.Element) {
	x, y, acc, err := cvutil.FindElementGenericWithoutRetry(canvas, lobbyScrollBarPath, mode)
	fmt.Printf("found lobby scroll at %v %v with accuracy: %v\n", x, y, acc)
	if err == nil {
		if mode < 2 {
			if acc < lobbyScrollAcc {
				util.Click(canvas, x, y)
			}
		} else {
			if acc > lobbyScrollAcc {
				util.Click(canvas, x, y)
			}
		}
	}
}

func findJoinButton(canvas *rod.Element) (int, int, error) {
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

func scrollLobby(canvas *rod.Element) {
	img, _ := cvutil.ElemToMat(canvas)
	scroll, _ := cvutil.LoadAssetToMat(lobbyScrollSpecialPath)
	x, y, acc := cvutil.GetMatchLocation(img, scroll, mode)
	fmt.Println(acc)
	if mode < 2 {
		if acc > lobbyScrollSpecialAcc {
			return
		}
	} else {
		if acc < lobbyScrollSpecialAcc {
			return
		}
	}

	shape, err := canvas.Shape()
	if err != nil {
		return
	}

	x += int(shape.Box().X)
	y += int(shape.Box().Y)

	tsz := scroll.Size()
	tx := tsz[0]
	ty := tsz[0]

	x += tx
	y += ty

	m := canvas.Page().Mouse
	m.Move(float64(x), float64(y), 1)
	for i := 0; i < 15; i++ {
		m.Click(proto.InputMouseButtonLeft)
	}
	fmt.Printf("found scroll bar at %v %v with acc %v", x, y, acc)
}
