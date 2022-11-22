package util

import (
	"github.com/shirou/gopsutil/process"
	"net/http"
	"encoding/json"
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

func GrabProcessByName(process_name string) *process.Process {
	procs, _ := process.Processes()
	for _, proc := range procs {
		name, _ := proc.Name()
		// _ since none of these should error...
		if (name == process_name) {
			return proc
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