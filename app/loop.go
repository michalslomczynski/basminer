package app

import (
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/michalslomczynski/bas-opencv/cvutil"
	"github.com/pkg/errors"
)

const (
	resetAfter = 10
)

func MainLopp() {
	canvas, err := Setup()
	if err != nil {
		log.Fatal(time.Now(), err)
	}

	for {
		err := Loop(canvas)
		if err != nil {
			log.Println(time.Now(), err)
			for {
				canvas, err = Restart(canvas)
				if err == nil {
					break
				}
			}
		}
	}
}

func Loop(canvas *rod.Element) error {
	state := &State{}
	assets, err := loadStateAssets(canvas)
	if err != nil {
		return err
	}

	retry := 0
	for {
		img, err := cvutil.ElemToMat(canvas)
		if err != nil {
			return err
		}

		StateUpdate(img, assets, state)

		if state.lobby {
			fmt.Println("found lobby state")
			err := SelectRoomAlt(canvas)
			if err != nil {
				log.Println(time.Now(), err)
			}
			retry = 0
			continue
		} else if state.round {
			fmt.Println("found round state")
			err := PlayRound(canvas)
			if err != nil {
				log.Println(time.Now(), err)
			}
			retry = 0
			continue
		} else if state.menu {
			fmt.Println("found menu state")
			err := GoToLobby(canvas)
			if err != nil {
				log.Println(time.Now(), err)
			}
			retry = 0
			continue
		} else {
			findAndClickWithoutRetry(canvas, joinFailedOkButtonPath, lobbyFailedOkButtonAcc)
			retry++
		}

		if retry > resetAfter {
			return errors.New("app stuck in main loop, aborting...")
		}
	}
}
