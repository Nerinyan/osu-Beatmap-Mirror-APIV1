package src

import (
	"encoding/json"
	"errors"
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
	go func() { //Loved
		for {
			time.Sleep(time.Second * 60)
			if err := getUpdatedMapLoved(); err != nil {
				ConsoleLogger.WarningConsolelog("Warning", err.Error())
				continue
			}
			if Settings.Config.Logger.UpdateSheduler {
				ConsoleLogger.UpdateLConsolelog("Update", "LOVED "+Settings.Config.Osu.BeatmapUpdate.UpdatedDesc.Id)
			}
		}
	}()
	go func() { //Qualified
		for {
			time.Sleep(time.Second * 60)
			if err := getUpdatedMapQualified(); err != nil {
				ConsoleLogger.WarningConsolelog("Warning", err.Error())
				continue
			}
			if Settings.Config.Logger.UpdateSheduler {
				ConsoleLogger.UpdateLConsolelog("Update", "QUALIFIED "+Settings.Config.Osu.BeatmapUpdate.UpdatedDesc.Id)
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
	go func() { //Update Graveyard asc limit 50
		for {
			time.Sleep(time.Minute)
			if err := getGraveyardMap(); err != nil {
				ConsoleLogger.WarningConsolelog("Warning", err.Error())
				continue
			}
			if Settings.Config.Logger.UpdateSheduler {
				ConsoleLogger.UpdateLConsolelog("Update", "GRAVEYARD "+Settings.Config.Osu.BeatmapUpdate.UpdatedAsc.Id)
			}
		}
	}()
}

func ManualUpdateBeatmapSet(id int) (err error) {
	url := fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d", id)

	var data osu.BeatmapSetsIN
	if err = stdGETBancho(url, &data); err != nil {
		return
	}

	updateMapset(&data)
	return
}

func getUpdatedMapRanked() (err error) {
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=ranked"

	var data osu.BeatmapsetsSearch
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}

	return
}

func getUpdatedMapLoved() (err error) {
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=loved"

	var data osu.BeatmapsetsSearch
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	return
}

func getUpdatedMapQualified() (err error) {
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=qualified"

	var data osu.BeatmapsetsSearch
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	return
}

func getUpdatedMapDesc() (err error) {
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_desc&s=any"

	var data osu.BeatmapsetsSearch

	if err = stdGETBancho(url, &data); err != nil {
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
	url := ""
	lu := &Settings.Config.Osu.BeatmapUpdate.UpdatedAsc.LastUpdate
	id := &Settings.Config.Osu.BeatmapUpdate.UpdatedAsc.Id
	if *lu+*id != "" {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=any&cursor%5Blast_update%5D=" + *lu + "&cursor%5B_id%5D=" + *id
	} else {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=any"
	}

	var data osu.BeatmapsetsSearch

	err = stdGETBancho(url, &data)
	if err != nil {
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
	*lu = *data.Cursor.LastUpdate
	*id = *data.Cursor.Id
	return
}

func getGraveyardMap() (err error) {

	url := ""
	lu := &Settings.Config.Osu.BeatmapUpdate.GraveyardAsc.LastUpdate
	id := &Settings.Config.Osu.BeatmapUpdate.GraveyardAsc.Id
	if *lu+*id != "" {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=graveyard&cursor%5Blast_update%5D=" + *lu + "&cursor%5B_id%5D=" + *id
	} else {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=graveyard"
	}

	var data osu.BeatmapsetsSearch

	err = stdGETBancho(url, &data)
	if err != nil {
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
	*lu = *data.Cursor.LastUpdate
	*id = *data.Cursor.Id
	return
}

func stdGETBancho(url string, str interface{}) (err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", Settings.Config.Osu.Token.TokenType+" "+Settings.Config.Osu.Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		res.Body.Close()
		apicountAdd()
	}()
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	return json.Unmarshal(body, &str)

}
