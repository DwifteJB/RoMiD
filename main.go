package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/DwifteJB/rblx-richpresence/util"
	"github.com/getlantern/systray"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/hugolgst/rich-go/client"
)

var current_placeId string
var useProfile = false
var profileDetails *util.RBLXPlayer
var profilePic string

type Settings struct {
	ShowProfile bool `json:"ShowProfile"`
	ClientId string `json:"ClientId"`
	RobloxId string `json:"RobloxId"`
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
	exist, err := util.DoesPathExist(c.homeDir)
	if err != nil {
		return err
	}
	if !exist {
		if err := os.Mkdir(c.homeDir, os.ModePerm); err != nil {
			return err
		}
	}

	cnf := path.Join(c.homeDir, c.configFile)
	exist, err = util.DoesPathExist(cnf)
	if err != nil {
		return err
	}
	if !exist {
		defaultSettings := Settings{
			ShowProfile: true,
			ClientId: "1044653106690015333",
			RobloxId: "156",
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



func onReady() {
	if err:= Config.Initalise(); err != nil {
		fmt.Print(err.Error())
	}
	if (Config.config.ShowProfile == true) {
		useProfile = true
		profileDetails = util.GetUserDetails(Config.config.RobloxId)
		if profileDetails == nil {
			fmt.Println("Couldn't get your roblox details. defaulting to off.")
			useProfile = false
		}
		profilePicDB := util.GetUserIcon(Config.config.RobloxId)
		if profilePicDB == nil || len(profilePicDB.Data) == 0 {
			fmt.Println("Couldn't get your roblox details. defaulting to off.")
			useProfile = false
		} else {
			profilePic = profilePicDB.Data[0].ImageURL
		}

		
	}
	userDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}
    systray.SetIcon(util.GetIcon("./src/icon.ico"))
    systray.SetTitle("Roblox Rich Presence")
    systray.SetTooltip("Roblox Rich Presence")

	name := systray.AddMenuItem("RblxPresence","RblxPresence")
	name.Disable()
	name.SetIcon(util.GetIcon("./src/icon.ico"))

	if useProfile == true {
		profile := systray.AddMenuItem("Using account "+profileDetails.Name,"Using account "+profileDetails.Name)
		profile.Disable()
	}
	systray.AddSeparator()

	connected := systray.AddMenuItem("Not connected to any game...", "Not connected...")
	connected.Disable()
	
	runOnOSBoot := systray.AddMenuItemCheckbox("Run on boot","Run on boot",false)

	exists, err := util.DoesPathExist(userDir + `\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\robloxRichPresence.lnk`)

	if err != nil {
		return
	}
	if exists {
		runOnOSBoot.Check()
	}
	
	openNotepad := systray.AddMenuItem("Open settings", "Open settings")

	systray.AddSeparator()

	quitMenu := systray.AddMenuItem("Close", "Quit the whole app")

	go func() {
		for {
			select {
			case <-quitMenu.ClickedCh:
				systray.Quit()
				return
			case <-openNotepad.ClickedCh:
				cmd := exec.Command("C:\\Windows\\system32\\notepad.exe", path.Join(Config.homeDir, Config.configFile))
				cmd.Run()
				return
			case <-runOnOSBoot.ClickedCh:
				if runOnOSBoot.Checked() == false {
					exists2, err := util.DoesPathExist(userDir + `\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\robloxRichPresence.lnk`)
					if err != nil {
						fmt.Println(err)
						return
					}
					if exists2 {
						os.Remove(userDir + `\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\robloxRichPresence.lnk`)
					}
					ex, err := os.Executable()

					fmt.Println(userDir)
					ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
					oleShellObject, err := oleutil.CreateObject("WScript.Shell")
					if err != nil {
						fmt.Println(err)
						return
					}
					defer oleShellObject.Release()
					wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
					if err != nil {
						fmt.Println(err)
						return
					}
					defer wshell.Release()
					
					cs, err := oleutil.CallMethod(wshell, "CreateShortcut", userDir + `\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\robloxRichPresence.lnk`)
					if err != nil {
						fmt.Println(err)
						return
					}
					idispatch := cs.ToIDispatch()
					_, err = oleutil.PutProperty(idispatch, "TargetPath", ex)
					if err != nil {
						fmt.Println(err)
						return
					}
		
					_, err = oleutil.PutProperty(idispatch, "Description", "Auto-run for Rblx-RichPresence")
					if err != nil {
						fmt.Println(err)
						return
					}
					_, err = oleutil.CallMethod(idispatch, "Save")
					if err != nil {
						fmt.Println(err)
						return
					}
					runOnOSBoot.Check()
					return
				} else {
					exists2, err := util.DoesPathExist(userDir + `\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\robloxRichPresence.lnk`)
					if err != nil {
						fmt.Println(err)
						return
					}
					if exists2 {
						os.Remove(userDir + `\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\robloxRichPresence.lnk`)
					}
					runOnOSBoot.Uncheck()
					return
				}
			}
		}
	}()
	errclient := client.Login(Config.config.ClientId)
	if errclient != nil {
		println("Couldn't start presence!")
		systray.SetTitle("Couldn't start presence!")
		systray.SetTooltip("Couldn't start presence!")
		name.SetTooltip("Couldn't start presence!")
	}
	fmt.Print("Presence is ready!\n")
	for {
		UpdatePresence(connected)
		time.Sleep(time.Second * 7)
	}
}

func onExit() {
    client.Logout()
}

func UpdatePresence(connected *systray.MenuItem) {
	placeId := util.GrabRobloxProcess()
	if placeId == "nil" {
		current_placeId = ""
		connected.SetTooltip("Not connected to any game...")
		connected.SetTitle("Not connected to any game...")
	} else if placeId != current_placeId {
		now := time.Now()
		place := util.GetPlaceInfoByPlaceId(placeId)
		placeIcon := util.GetIconByPlaceId(placeId)
		println(place)
		if place == nil || placeIcon == nil {
			fmt.Println("Couldn't get the games details..")
			return
		}
		connected.SetTooltip("Connected to " + place.Name + " by " + place.Creator.Name)
		connected.SetTitle("Connected to " + place.Name + " by " + place.Creator.Name)
		Activity := client.Activity{}
		if useProfile == true {
			Activity = client.Activity{
				State: "by " + place.Creator.Name,
				Details: "Playing "+ place.Name,
				LargeImage: placeIcon.Data[0].ImageURL,
				LargeText: "RBLX Presence 1.0 | Created by Dwifte",
				SmallImage: profilePic,
				SmallText: "Logged in as "+profileDetails.Name,
				Buttons: []*client.Button {
					{
						Label: "Play this game",
						Url: "https://www.roblox.com/games/" + placeId + "/-",
					},
				},
				Timestamps: &client.Timestamps {
					Start: &now,
				},
			}
		} else {
			Activity = client.Activity{
				State: "by " + place.Creator.Name,
				Details: "Playing "+ place.Name,
				LargeImage: placeIcon.Data[0].ImageURL,
				LargeText: "RBLX Presence 1.0 | Created by Dwifte",
				Buttons: []*client.Button {
					{
						Label: "Play this game",
						Url: "https://www.roblox.com/games/" + placeId + "/-",
					},
				},
				Timestamps: &client.Timestamps {
					Start: &now,
				},
			}
		}
		client.SetActivity(Activity)
		println("Connected to " + place.Name + " by " + place.Creator.Name)
		current_placeId = placeId
	}
}

func main() {
    systray.Run(onReady, onExit)
}