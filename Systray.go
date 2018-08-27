package main

import (
	"log"
	"os"
	"path"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

var (
	exePath                               string
	mInfo, mEnable, mReload, mHelp, mQuit *systray.MenuItem
)

func main() {
	var err error
	exePath, err = os.Executable()
	if err != nil {
		log.Panic(err)
	}
	exePath = path.Dir(exePath)
	systray.Run(onReady, onExit)
}

func onExit() {
	enable(false)
}

func onReady() {
	mInfo = systray.AddMenuItem(versionStr, authorStr)
	mInfo.Disable()
	systray.AddSeparator()
	mEnable = systray.AddMenuItem("Pause", "Pause Gosture")
	mReload = systray.AddMenuItem("Reload Config", "Reload configuraton file")
	mHelp = systray.AddMenuItem("Help", "Open user manual")
	systray.AddSeparator()
	mQuit = systray.AddMenuItem("Quit", "Quit Gosture")
	systray.SetIcon(iconON)
	systray.SetTitle(versionStr)
	systray.SetTooltip(authorStr)
	go gosture()
	for {
		select {
		case <-mInfo.ClickedCh:
		case <-mEnable.ClickedCh:
			enable(!ena)
			onToggle()
		case <-mReload.ClickedCh:
			reload()
		case <-mHelp.ClickedCh:
			open.Run(path.Join(exePath, "Gosture_Help.htm"))
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func onToggle() {
	if ena {
		systray.SetIcon(iconON)
		mEnable.SetTitle("Pause")
		mEnable.SetTooltip("Pause Gosture")
	} else {
		systray.SetIcon(iconOFF)
		mEnable.SetTitle("Resume")
		mEnable.SetTooltip("Resume Gosture")
	}
	if rdy {
		mEnable.Enable()
	} else {
		mEnable.Disable()
	}
}
