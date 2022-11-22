package util

import (
	"github.com/shirou/gopsutil/process"
	"net/http"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"fmt"
	"regexp"
)

type RBLXReturnData struct {
	Data []struct {
		TargetID int64  `json:"targetId"`
		State    string `json:"state"`
		ImageURL string `json:"imageUrl"`
	} `json:"data"`
}

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

func GrabRobloxProcess() *process.Process {
	procs, _ := process.Processes()
	for _, proc := range procs {
		name, _ := proc.Name()
		// _ since none of these should error...
		if (name == "RobloxPlayerBeta.exe") {
			cmdLine, _ := proc.Cmdline()
			placePattern := regexp.MustCompile(`placeId=(\d+)`)
			placeMatch := placePattern.FindStringSubmatch(cmdLine)
			if len(placeMatch) != 0 {
				return proc
			}

		}
	}
	return nil
}

func GetPlaceInfoByPlaceId(placeId string) *MarketPlaceInfo {
	url := "https://api.roblox.com/marketplace/productinfo?assetId=" + placeId
	resp, err := http.Get(url)
	if err != nil {
		print("Got an error fetching info: ", err.Error(),"\n")
	}
	defer resp.Body.Close()
	var info *MarketPlaceInfo
	json.NewDecoder(resp.Body).Decode(&info)
	return info
}

func GetIconByPlaceId(placeId string) * RBLXReturnData {
	url := "https://thumbnails.roblox.com/v1/assets?returnPolicy=0&size=512x512&format=Png&isCircular=false&assetIds=" + placeId
	resp, err := http.Get(url)
	if err != nil {
		print("Got an error fetching info: ", err.Error(),"\n")
	}
	defer resp.Body.Close()
	var info *RBLXReturnData
	json.NewDecoder(resp.Body).Decode(&info)
	return info
}

func GetUserIcon(userId string) *RBLXReturnData {
	url := "https://thumbnails.roblox.com/v1/assets?returnPolicy=0&size=512x512&format=Png&isCircular=false&assetIds=" + userId
	resp, err := http.Get(url)
	if err != nil {
		print("Got an error fetching info: ", err.Error(),"\n")
	}
	defer resp.Body.Close()
	var info *RBLXReturnData
	json.NewDecoder(resp.Body).Decode(&info)
	return info
}
func GetIcon(s string) []byte {
    b, err := ioutil.ReadFile(s)
    if err != nil {
        fmt.Print(err)
    }
    return b
}


func DoesPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}