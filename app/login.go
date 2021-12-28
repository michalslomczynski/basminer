package app

import (
	"github.com/go-rod/rod"
	"github.com/go-vgo/robotgo"
	"github.com/michalslomczynski/bas-opencv/config"
	"github.com/michalslomczynski/bas-opencv/cvutil"
	"github.com/michalslomczynski/bas-opencv/util"
	"github.com/michalslomczynski/bas-opencv/web/wallet"
	cv "gocv.io/x/gocv"
	"time"
)

const (
	timeout = 60 * time.Second
	mode    = cv.TmSqdiffNormed
	// Paths
	prefix                   = "assets/"
	acceptButtonDisabledPath = prefix + "accept_disabled.png"
	acceptButtonEnabledPath  = prefix + "accept_enabled.png"
	usernamePath             = prefix + "enter_username.png"
	connectToWalletPath      = prefix + "connect_to_wallet.png"
	playButtonPath           = prefix + "play_button.png"
	// Accuracies for SqDiff
	acceptButtonDisabledAcc  = 0.01
	acceptButtonEnabledAcc   = 0.04
	usernameAcc              = 0.003
	connectToWalletButtonAcc = 0.05
)

func AcceptTerms(canvas *rod.Element) error {
	x, y, err := scrollUntilAcceptButtonEnabled(canvas)
	if err != nil {
		return err
	}

	util.Click(canvas, x, y)

	return nil
}

func Login(canvas *rod.Element) error {
	x, y, err := FindUsername(canvas)
	if err != nil {
		return err
	}
	util.Click(canvas, x, y)

	// TODO: replace with working native browser events - requires isTrusted property
	robotgo.TypeStr(config.Cfg.Username)

	err = wallet.SelectMetaMask(canvas.Page())
	if err != nil {
		return err
	}

	err = wallet.ConnectWithMetamask(canvas.Page())
	if err != nil {
		return err
	}

	err = wallet.SignTransaction(canvas.Page())
	if err != nil {
		return err
	}

	x, y, err = PlayButton(canvas)
	if err != nil {
		return err
	}

	util.Click(canvas, x, y)

	return nil
}

func LoginAfterRestart(canvas *rod.Element) error {
	// TODO: replace with working native browser events - requires isTrusted property
	robotgo.TypeStr(config.Cfg.Username)

	x, y, err := ConnectToWalletButton(canvas)
	if err != nil {
		return err
	}
	util.Click(canvas, x, y)

	err = wallet.SelectMetaMask(canvas.Page())
	if err != nil {
		return err
	}

	err = wallet.SignTransaction(canvas.Page())
	if err != nil {
		return err
	}

	x, y, err = PlayButton(canvas)
	if err != nil {
		return err
	}

	util.Click(canvas, x, y)

	return nil
}

func FindAcceptButtonDisabled(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, acceptButtonDisabledPath, mode, acceptButtonDisabledAcc, timeout)
}

func FindAcceptButtonEnabled(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, acceptButtonEnabledPath, mode, acceptButtonEnabledAcc, timeout)
}

func FindUsername(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, usernamePath, mode, usernameAcc, timeout)
}

func ConnectToWalletButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, connectToWalletPath, mode, connectToWalletButtonAcc, timeout)
}

func PlayButton(canvas *rod.Element) (int, int, error) {
	return cvutil.FindElementGeneric(canvas, playButtonPath, mode, battleModeButtonAcc, timeout)
}

func scrollUntilAcceptButtonEnabled(canvas *rod.Element) (int, int, error) {
	x, y, err := FindAcceptButtonDisabled(canvas)
	if err != nil {
		return 0, 0, err
	}

	// Blind guess offset to scrollable.
	oy := float64(-50)

	for {
		for i := 0; i < 50; i++ {
			canvas.Page().Mouse.Move(float64(x), float64(y)+oy, 1)
			canvas.Page().Mouse.Scroll(0, 10000000, 1)
		}

		_, _, acc, err := cvutil.FindElementGenericWithoutRetry(canvas, acceptButtonEnabledPath, mode)
		if err != nil {
			return 0, 0, err
		}

		if acc < acceptButtonEnabledAcc {
			return x, y, nil
		}
	}
}
