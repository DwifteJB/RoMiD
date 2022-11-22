package main

import (
	"encoding/json"
	//"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
	"github.com/DwifteJB/rblx-richpresence/util"
	"github.com/getlantern/systray"
	"github.com/hugolgst/rich-go/client"
)
type MarketPlaceInfo struct { 
	Name        string      `json:"Name"`
	Description string      `json:"Description"`
	Creator     struct {
		Id              int    `json:"Id"`
		Name            string `json:"Name"`
		CreatorType     string `json:"CreatorType"`
		CreatorTargetId int    `json:"CreatorTargetId"`
	} `json:"Creator"`
	IconImageAssetId       int64       `json:"IconImageAssetId"`
}
type Settings struct {
	ShowUsername bool `json:"ShowUsername"`
	ClientId string `json:"ClientId"`
}


var Config = func() *config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir, _ = os.Getwd()
	}

	homeDir = path.Join(homeDir, ".dwifte")
	return &config{
		homeDir:    homeDir,
		configFile: "roblox-rich-presence.json",
		config: Settings{},
	}
}()

type config struct {
	homeDir    string
	configFile string
	config Settings

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
		defaultSettings := Settings{
			ShowUsername: true,
			ClientId: "1044653106690015333",
		}
		data, _:= json.Marshal(defaultSettings)
		ioutil.WriteFile(cnf,data,0644)
	} else if data, err = os.ReadFile(cnf); err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &c.config)
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
	err := client.Login(Config.config.ClientId)
	if err != nil {
		println("Couldn't start presence!")
		systray.SetTitle("Couldn't start presence!")
		systray.SetTooltip("Couldn't start presence!")
		name.SetTooltip("Couldn't start presence!")
	}
	client.SetActivity(client.Activity {
		Details: "Waiting to join a game...",
		LargeImage: "https://github.com/DwifteJB.png",
		LargeText: "Playing Roblox!",
	})
	fmt.Print("Presence is ready!\n")
	for {
		time.Sleep(time.Second * 5)
	}
}

func onExit() {
    client.Logout()
}

