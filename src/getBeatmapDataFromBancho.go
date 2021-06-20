package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/nerina1241/osu-beatmap-mirror-api/ConsoleLogger"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
	"github.com/nerina1241/osu-beatmap-mirror-api/osu"
)

var api = struct {
	count int
	mutex sync.Mutex
}{}

func apicountAdd() {
	api.mutex.Lock()
	api.count++
	api.mutex.Unlock()
}
func apiCountReset() {
	api.mutex.Lock()
	api.count = 0
	api.mutex.Unlock()
}

func awaitApiCount() {
	for {
		if api.count < 60 {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func RunGetBeatmapDataASBancho() {
	checkUpdatable := Settings.Config.Osu.Token.UpdatedAt + Settings.Config.Osu.Token.ExpiresIn - time.Now().Unix()
	if checkUpdatable < 3600 {
		time.Sleep(time.Second * 10)
	}

	go func() {
		for {
			time.Sleep(time.Minute)

			if Maria.Ping() != nil {
				continue
			}
			apiCountReset()
			go Settings.Config.Save()
		}
	}()
	go func() { //desc
		for {
			time.Sleep(time.Second * 30)
			if err := getUpdatedMapDesc(); err != nil {
				ConsoleLogger.WarningConsolelog("Warning", err.Error())
				continue
			}
			if Settings.Config.Logger.UpdateSheduler {
				ConsoleLogger.UpdateLConsolelog("Update", "DESC "+Settings.Config.Osu.BeatmapUpdate.UpdatedDesc.Id)
			}
		}
	}()
	go func() { //Ranked
		for {
			time.Sleep(time.Second * 60)
			if err := getUpdatedMapRanked(); err != nil {
				ConsoleLogger.WarningConsolelog("Warning", err.Error())
				continue
			}
			if Settings.Config.Logger.UpdateSheduler {
				ConsoleLogger.UpdateLConsolelog("Update", "RANKED "+Settings.Config.Osu.BeatmapUpdate.UpdatedDesc.Id)
			}
		}
	}()
	go func() { //asc
		for {
			awaitApiCount()

			if err := getUpdatedMapAsc(); err != nil {
				ConsoleLogger.WarningConsolelog("Warning", err.Error())
				continue
			}
			if Settings.Config.Logger.UpdateSheduler {
				ConsoleLogger.UpdateLConsolelog("Update", "ASC "+Settings.Config.Osu.BeatmapUpdate.UpdatedAsc.Id)
			}
		}
	}()
}

func ManualUpdateBeatmapSet(id int) (err error) {
	url := fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d", id)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}
	req.Header.Add("Authorization", Settings.Config.Osu.Token.TokenType+" "+Settings.Config.Osu.Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}
	defer func() {
		res.Body.Close()
		apicountAdd()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}
	ms := string(body)
	if ms == "" || ms == "{\"error\":null}" || res.StatusCode != 200 {
		return
	}

	var v osu.BeatmapSets
	if err = json.Unmarshal([]byte(ms), &v); err != nil {
		fmt.Print(id, "error", err.Error())
		return
	}
	updateMapset(&v)
	return
}

func getUpdatedMapRanked() (err error) {
	//TODO 30sec

	//https://osu.ppy.sh/beatmapsets/search?sort=updated_desc&s=any&cursor%5Blast_update%5D=1621954136000&cursor%5B_id%5D=1473132
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=ranked"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}

	req.Header.Add("Authorization", Settings.Config.Osu.Token.TokenType+" "+Settings.Config.Osu.Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}
	defer func() {
		res.Body.Close()
		apicountAdd()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}

	var data osu.BeatmapsetsSearch
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}

	return
}

func getUpdatedMapDesc() (err error) {
	//TODO 30sec

	//https://osu.ppy.sh/beatmapsets/search?sort=updated_desc&s=any&cursor%5Blast_update%5D=1621954136000&cursor%5B_id%5D=1473132
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_desc&s=any"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}

	req.Header.Add("Authorization", Settings.Config.Osu.Token.TokenType+" "+Settings.Config.Osu.Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}
	defer func() {
		res.Body.Close()
		apicountAdd()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}

	var data osu.BeatmapsetsSearch
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	Settings.Config.Osu.BeatmapUpdate.UpdatedDesc.LastUpdate = *data.Cursor.LastUpdate
	Settings.Config.Osu.BeatmapUpdate.UpdatedDesc.Id = *data.Cursor.Id

	return
}

func getUpdatedMapAsc() (err error) {
	//TODO

	//      https://osu.ppy.sh/beatmapsets/search?sort=updated_desc&s=any&cursor%5Blast_update%5D=1621954136000&cursor%5B_id%5D=1473132
	//      https://osu.ppy.sh/beatmapsets/search?sort=updated_desc&s=any&cursor%5Blast_update%5D=1622554856000&cursor%5B_id%5D=1477878
	url := ""
	lu := &Settings.Config.Osu.BeatmapUpdate.UpdatedAsc.LastUpdate
	id := &Settings.Config.Osu.BeatmapUpdate.UpdatedAsc.Id
	if *lu+*id != "" {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=any&cursor%5Blast_update%5D=" + *lu + "&cursor%5B_id%5D=" + *id
	} else {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=any"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}

	req.Header.Add("Authorization", Settings.Config.Osu.Token.TokenType+" "+Settings.Config.Osu.Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}
	defer func() {
		res.Body.Close()
		apicountAdd()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}

	var data osu.BeatmapsetsSearch
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}
	if data.Cursor == nil {
		*lu = ""
		*id = ""
		return
	}

	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	//fmt.Println(*lu, *id, *data.Cursor , data.Beatmapsets == nil)
	*lu = *data.Cursor.LastUpdate
	*id = *data.Cursor.Id
	return
}
