package main

import (
    "fmt"
    "io/ioutil"
	"os"
	"path"
	"errors"
	//"encoding/json"
	"github.com/hugolgst/rich-go/client"
    "github.com/getlantern/systray"
)

type Settings struct {
	ShowUsername bool `json:"ShowUsername"`
	ClientId int `json:"ClientId"`
}

type MarketPlaceInfo struct { // https://mholt.github.io/json-to-go/
	Name        string      `json:"Name"`
	Description string      `json:"Description"`
	Creator     struct {
		ID              int    `json:"Id"`
		Name            string `json:"Name"`
		CreatorType     string `json:"CreatorType"`
		CreatorTargetID int    `json:"CreatorTargetId"`
	} `json:"Creator"`
	IconImageAssetID       int64       `json:"IconImageAssetId"`
}

var Config = func() *config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir, _ = os.Getwd()
	}

	homeDir = path.Join(homeDir, ".dwifte")
	return &config{
		homeDir:    homeDir,
		configFile: "roblox-rich-presence.json"}
}()

type config struct {
	homeDir    string
	configFile string

}

func (c *config) Initalise() error {
	var data []byte
	exist, err := doesPathExist(c.homeDir)
	if err != nil {
		return err
	}
	if !exist {
		if err := os.Mkdir(c.homeDir, os.ModePerm); err != nil {
			return err
		}
	}

	cnf := path.Join(c.homeDir, c.configFile)
	exist, err = doesPathExist(cnf)
	if err != nil {
		return err
	}
	if !exist {
		f, err := os.Create(cnf)
		if err != nil {
			return err
		}
		f.Close()
	} else if data, err = os.ReadFile(cnf); err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return nil
}

func main() {
    systray.Run(onReady, onExit)
}

func onReady() {
    systray.SetIcon(getIcon("./src/icon.ico"))
    systray.SetTitle("Roblox Rich Presence")
    systray.SetTooltip("Roblox Rich Presence")

	name := systray.AddMenuItem("RblxPresence","RblxPresence")
	name.Disable()
	name.SetIcon(getIcon("./src/icon.ico"))

	systray.AddSeparator()

	connected := systray.AddMenuItem("Not connected to any game...", "Not connected...")
	connected.Disable()

	systray.AddSeparator()

	quitMenu := systray.AddMenuItem("Close", "Quit the whole app")
	if err:= Config.Initalise(); err != nil {
		fmt.Print(err.Error())
	}
	go func() {
		for {
			select {
			case <-quitMenu.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
	fmt.Print("Presence is ready!\n")
}

func onExit() {
    // Cleaning stuff here.
}

func getIcon(s string) []byte {
    b, err := ioutil.ReadFile(s)
    if err != nil {
        fmt.Print(err)
    }
    return b
}


func doesPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}