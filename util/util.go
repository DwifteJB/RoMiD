package util

import (
	"github.com/shirou/gopsutil/v3/process"
	"net/http"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"fmt"
	"regexp"
	"time"
)
type RBLXPlayer struct {
	Description            string      `json:"description"`
	Created                time.Time   `json:"created"`
	IsBanned               bool        `json:"isBanned"`
	ExternalAppDisplayName interface{} `json:"externalAppDisplayName"`
	HasVerifiedBadge       bool        `json:"hasVerifiedBadge"`
	ID                     int         `json:"id"`
	Name                   string      `json:"name"`
	DisplayName            string      `json:"displayName"`
}
type RBLXError struct {
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}
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

func GrabRobloxProcess() string {
	procs, _ := process.Processes()
	for _, proc := range procs {
		name, _ := proc.Name()
		if (name == "RobloxPlayerBeta.exe") {
			cmdLine, _ := proc.Cmdline()
			placePattern := regexp.MustCompile(`placeId=(\d+)`)
			if len(placePattern.FindStringSubmatch(cmdLine)) != 0 {
				return placePattern.FindStringSubmatch(cmdLine)[1]
			}
		}
	}
	return "nil"
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

func GetIconByPlaceId(placeId string) *RBLXReturnData {
	url := "https://thumbnails.roblox.com/v1/places/gameicons?returnPolicy=0&size=512x512&format=Png&isCircular=false&placeIds=" + placeId
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
	url := "https://thumbnails.roblox.com/v1/users/avatar-headshot?size=420x420&format=Png&isCircular=false&userIds=" + userId
	resp, err := http.Get(url)
	if err != nil {
		print("Got an error fetching info: ", err.Error(),"\n")
	}
	defer resp.Body.Close()
	var info *RBLXReturnData

	json.NewDecoder(resp.Body).Decode(&info)

	return info
}

func GetUserDetails(userId string) *RBLXPlayer{
	
	url := "https://users.roblox.com/v1/users/" + userId
	resp, err := http.Get(url)
	if err != nil {
		print("Got an error fetching info: ", err.Error(),"\n")
	}
	defer resp.Body.Close()
	var info *RBLXPlayer
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