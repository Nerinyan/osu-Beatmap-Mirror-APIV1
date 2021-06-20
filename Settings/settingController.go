package Settings

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/nerina1241/osu-beatmap-mirror-api/ConsoleLogger"
)

type config struct {
	Port      string `json:"port"`
	TargetDir string `json:"targetDir"`
	Logger    struct {
		UpdateSheduler     bool `json:"updatesSheduler"`
		DownloadBeatmap    bool `json:"downloadBeatmap"`
		UpdateBeatmap      bool `json:"updateBeatmap"`
		ShowFavouriteCount struct {
			ALL    bool `json:"all"`
			Over70 bool `json:"over70"`
		} `json:"showFavouriteCount"`
	} `json:"logger"`
	AutoDownload70FavOver bool   `json:"autoDownload70FavOver"`
	Key                   string `json:"Key"`
	Sql                   struct {
		Id     string `json:"id"`
		Passwd string `json:"passwd"`
		Url    string `json:"url"`
	} `json:"sql"`
	Osu struct {
		Username string `json:"username"`
		Passwd   string `json:"passwd"`
		Token    struct {
			TokenType    string `json:"token_type"`
			ExpiresIn    int64  `json:"expires_in"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			UpdatedAt    int64  `json:"updatedAt"`
		} `json:"token"`
		BeatmapUpdate struct {
			UpdatedAsc struct {
				LastUpdate string `json:"last_update"`
				Id         string `json:"_id"`
			} `json:"updated_asc"`
			UpdatedDesc struct {
				LastUpdate string `json:"last_update"`
				Id         string `json:"_id"`
			} `json:"updated_desc"`
		} `json:"beatmapUpdate"`
	} `json:"osu"`
}

var Config config

func LoadSetting() {
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		ConsoleLogger.WarningConsolelog("CONFIG", "i can't find config, so i make a new for you")
		Config.Save()
		os.Exit(3)
	}
	err = json.Unmarshal(b, &Config)
	if err != nil {
		ConsoleLogger.WarningConsolelog("CONFIG", "idk, your config file has something wrong. details: "+err.Error())
		os.Exit(3)
	}
}

func (v *config) Save() {
	file, _ := json.MarshalIndent(v, "", "  ")
	_ = ioutil.WriteFile("config.json", file, 0755)
}
