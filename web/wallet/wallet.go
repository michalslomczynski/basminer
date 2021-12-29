package wallet

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/michalslomczynski/bas-opencv/browser"
	"github.com/michalslomczynski/bas-opencv/config"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const (
	extensionUrl           = "chrome-extension://nojmdacjahoombflmhcebljpcpjfogei/home.html#initialize/welcome"
	extensionPageTimeout   = 40
	extensionPageTitle     = "MetaMask"
	selectorsTimeout       = 2 * time.Second
	transactionPageTitle   = "MetaMask Notification"
	transactionPageTimeout = 20
	// Binance network connection data
	networkName      = "Smart Chain"
	RPCURL           = "https://bsc-dataseed.binance.org/"
	chainId          = "56"
	currencySymbol   = "BNB"
	blockExplorerURL = " https://bscscan.com"
	// Selectors
	getStartedSelector                     = "#app-content > div > div.main-container-wrapper > div > div > div > button"
	importWalletSelector                   = "#app-content > div > div.main-container-wrapper > div > div > div.select-action__wrapper > div > div.select-action__select-buttons > div:nth-child(1) > button"
	agreeSelector                          = "#app-content > div > div.main-container-wrapper > div > div > div > div.metametrics-opt-in__footer > div.page-container__footer > footer > button.button.btn--rounded.btn-primary.page-container__footer-button"
	secretPhaseSelector                    = "#app-content > div > div.main-container-wrapper > div > div > form > div.first-time-flow__textarea-wrapper > div.MuiFormControl-root.MuiTextField-root.first-time-flow__textarea.first-time-flow__seedphrase > div > input"
	newPasswdSelector                      = "#password"
	newPasswdConfirmSelector               = "#confirm-password"
	termsTickBoxSelector                   = "#app-content > div > div.main-container-wrapper > div > div > form > div.first-time-flow__checkbox-container > div"
	importButtonSelector                   = "#app-content > div > div.main-container-wrapper > div > div > form > button"
	allDoneSelector                        = "#app-content > div > div.main-container-wrapper > div > div > button"
	whatIsNewCloseSelector                 = "#popover-content > div > div > section > header > div > button"
	selectMetaMaskWalletSelector           = "body > aside > section > ul > li > button > span"
	signWalletTransactionSelector          = "#app-content > div > div.main-container-wrapper > div > div.permissions-connect-choose-account > div.permissions-connect-choose-account__footer-container > div.permissions-connect-choose-account__bottom-buttons > button.button.btn--rounded.btn-primary"
	signWalletTransactionConfirmSelector   = "#app-content > div > div.main-container-wrapper > div > div.page-container.permission-approval-container > div.permission-approval-container__footers > div.page-container__footer > footer > button.button.btn--rounded.btn-primary.page-container__footer-button"
	authorizeTransactionSignButtonSelector = "#app-content > div > div.main-container-wrapper > div > div.request-signature__footer > button.button.btn--rounded.btn-primary.btn--large.request-signature__footer__sign-button"
	dropDownNetworkSelector                = "#app-content > div > div.app-header.app-header--back-drop > div > div.app-header__account-menu-container > div.app-header__network-component-wrapper > div > div.chip__right-icon > i"
	addNetworkButtonSelector               = "#app-content > div > div.menu-droppo-container.network-droppo > div > button"
	networkNameInputSelecotr               = "#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.networks-tab__content > div > div.networks-tab__add-network-form-body > div:nth-child(1) > label > input"
	RPCURLInputSelector                    = "#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.networks-tab__content > div > div.networks-tab__add-network-form-body > div:nth-child(2) > label > input"
	chainIdInputSelector                   = "#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.networks-tab__content > div > div.networks-tab__add-network-form-body > div:nth-child(3) > label > input"
	currencySymbolInputSelector            = "#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.networks-tab__content > div > div.networks-tab__add-network-form-body > div:nth-child(4) > label > input"
	blockExplorerInputSelector             = "#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.networks-tab__content > div > div.networks-tab__add-network-form-body > div:nth-child(5) > label > input"
	saveNewNetworkButtonSelector           = "#app-content > div > div.main-container-wrapper > div > div.settings-page__content > div.settings-page__content__modules > div > div.networks-tab__content > div > div.networks-tab__add-network-form-footer > button.button.btn--rounded.btn-primary"
)

// SignInToWallet logins to MetaMask wallet plugin in browser.
func SignInToWallet(b *rod.Browser) (*rod.Page, error) {
	// Force browser to focus
	b.MustPage("")

	page := waitForMetaMaskPageWithTimeout(b, extensionPageTimeout)
	if page == nil {
		err := rod.Try(func() {
			page = b.MustPage(extensionUrl).MustWaitLoad()
		})
		if err != nil {
			return nil, errors.Wrapf(err, "could not open metamask page")
		}
	}

	err := loginToWallet(page)
	if err != nil {
		return nil, errors.Wrapf(err, "could not login to wallet")
	}

	err = ConnectBinanceNetwork(page)
	if err != nil {
		return nil, errors.Wrapf(err, "could not add binance network to wallet")
	}

	return page, nil
}

func waitForMetaMaskPageWithTimeout(b *rod.Browser, timeout time.Duration) *rod.Page {
	c := make(chan *rod.Page, 1)

	targetCondition := func(t *proto.TargetTargetInfo) bool {
		return t.Title == extensionPageTitle
	}

	go func() {
		for {
			target := browser.GetTarget(b, targetCondition)
			if target == nil {
				continue
			}
			if target.Type == proto.TargetTargetInfoTypePage {
				page, err := b.PageFromTarget(target.TargetID)
				if err != nil {
					c <- nil
				}
				c <- page
			}
		}
	}()

	select {
	case page := <-c:
		return page
	case <-time.After(timeout * time.Second):
		return nil
	}
}

func loginToWallet(page *rod.Page) error {
	err := rod.Try(func() {
		page.Timeout(selectorsTimeout)
		page.MustElement(getStartedSelector).MustClick()
		page.MustElement(importWalletSelector).MustClick()
		page.MustElement(agreeSelector).MustClick()
		page.MustElement(secretPhaseSelector).MustInput(config.Cfg.WalletPassphrase)
		page.MustElement(newPasswdSelector).MustInput(config.Cfg.WalletPassword)
		page.MustElement(newPasswdConfirmSelector).MustInput(config.Cfg.WalletPassword)
		page.MustElement(termsTickBoxSelector).MustClick()
		page.MustElement(importButtonSelector).MustClick()
		page.MustElement(allDoneSelector).MustClick()
		page.MustElement(whatIsNewCloseSelector).MustClick()
	})
	if err != nil {
		return err
	}
	return nil
}

func ConnectWithMetamask(page *rod.Page) error {
	walletPage, err := findMetaMaskNotificationPage(page)
	if err != nil {
		return err
	}

	err = rod.Try(func() {
		walletPage.Timeout(selectorsTimeout)
		walletPage.MustElement(signWalletTransactionSelector).MustClick()
		walletPage.MustElement(signWalletTransactionConfirmSelector).MustClick()
	})
	if err != nil {
		return err
	}

	return nil
}

func SelectMetaMask(page *rod.Page) error {
	el, err := page.Element(selectMetaMaskWalletSelector)
	if err != nil {
		return err
	}
	el.Click(proto.InputMouseButtonLeft)

	return nil
}

func SignTransaction(page *rod.Page) error {
	pageSearchOpt := func(t *proto.TargetTargetInfo) bool {
		return strings.Contains(t.URL, "transaction")
	}

	walletPage, err := findMetaMaskNotificationPage(page, pageSearchOpt)
	if err != nil {
		return err
	}

	el, err := walletPage.Element(authorizeTransactionSignButtonSelector)
	if err != nil {
		return errors.New("could not find sign in button")
	}
	el.Click(proto.InputMouseButtonLeft)

	return nil
}

func findMetaMaskNotificationPage(page *rod.Page, opts ...browser.TargetOpts) (*rod.Page, error) {
	fmt.Println("Looking for notification page")

	targetTitleCond := func(t *proto.TargetTargetInfo) bool {
		return t.Title == transactionPageTitle
	}
	opts = append(opts, targetTitleCond)

	page.MustWaitLoad()
	b := page.Browser()

	var target *proto.TargetTargetInfo
	for i := 0; i < transactionPageTimeout; i++ {
		target = browser.GetTarget(b, opts...)
		if target != nil {
			break
		}
		time.Sleep(time.Second)
	}
	fmt.Println("DEBUG: ", target)
	if target == nil {
		return nil, errors.New(fmt.Sprintf("could not find page %+v with any target", transactionPageTitle))
	}

	walletPage, err := b.PageFromTarget(target.TargetID)
	if err != nil {
		return nil, err
	}
	fmt.Println("DEBUG: found target id", walletPage.TargetID)

	return walletPage, nil
}

func ConnectBinanceNetwork(page *rod.Page) error {
	err := rod.Try(func() {
		page.Timeout(selectorsTimeout)
		page.MustElement(dropDownNetworkSelector).MustClick()
		page.MustElement(addNetworkButtonSelector).MustClick().MustWaitLoad()
		page.MustElement(networkNameInputSelecotr).MustInput(networkName)
		page.MustElement(RPCURLInputSelector).MustInput(RPCURL)
		page.MustElement(chainIdInputSelector).MustInput(chainId)
		page.MustElement(currencySymbolInputSelector).MustInput(currencySymbol)
		page.MustElement(saveNewNetworkButtonSelector).MustClick().MustWaitLoad()
	})
	if err != nil {
		return err
	}
	return nil
}
