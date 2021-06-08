package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
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

func RunGetBeatmapDataASBancho() {
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
				fmt.Println(err)
				continue
			}
			if Settings.Config.Logger.UpdateSheduler {
				fmt.Println("[U]", "DESC", Settings.Config.Osu.BeatmapUpdate.UpdatedDesc.Id)
			}
		}
	}()
	go func() { //Ranked
		for {
			time.Sleep(time.Second * 60)
			if err := getUpdatedMapRanked(); err != nil {
				fmt.Println(err)
				continue
			}
			if Settings.Config.Logger.UpdateSheduler {
				fmt.Println("[U]", "RANKED", Settings.Config.Osu.BeatmapUpdate.UpdatedDesc.Id)
			}
		}
	}()
	go func() { //asc
		for {
			awaitApiCount()

			if err := getUpdatedMapAsc(); err != nil {
				fmt.Println(err)
				continue
			}
			if Settings.Config.Logger.UpdateSheduler {
				fmt.Println("[U]", "ASC", Settings.Config.Osu.BeatmapUpdate.UpdatedAsc.Id)
			}
		}
	}()
}

func awaitApiCount() {
	for {
		if api.count < 60 {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func ManualUpdateBeatmapSet(id string) (err error) {
	ms := getBeatmapSets(id)
	if ms == "" || ms == "{\"error\":null}" {
		return errors.New("")
	}

	var v map[string]interface{}
	if err = json.Unmarshal([]byte(ms), &v); err != nil {
		return
	}
	updateMap(v)
	return
}

func getBeatmapSets(id string) string {
	url := "https://osu.ppy.sh/api/v2/beatmapsets/" + id
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Authorization", Settings.Config.Osu.Token.TokenType+" "+Settings.Config.Osu.Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer func() {
		res.Body.Close()
		apicountAdd()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}

func getUpdatedMapDesc() (err error) {
	//TODO 30sec

	//https://osu.ppy.sh/beatmapsets/search?sort=updated_desc&s=any&cursor%5Blast_update%5D=1621954136000&cursor%5B_id%5D=1473132
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_desc&s=any"

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data map[string]interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	if err = updateSearchBeatmaps(data); err != nil {
		return
	}
	c := data["cursor"].(map[string]interface{})
	Settings.Config.Osu.BeatmapUpdate.UpdatedDesc.LastUpdate = c["last_update"].(string)
	Settings.Config.Osu.BeatmapUpdate.UpdatedDesc.Id = c["_id"].(string)
	return
}
func getUpdatedMapRanked() (err error) {
	//TODO 30sec

	//https://osu.ppy.sh/beatmapsets/search?sort=updated_desc&s=any&cursor%5Blast_update%5D=1621954136000&cursor%5B_id%5D=1473132
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=ranked"

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data map[string]interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	if err = updateSearchBeatmaps(data); err != nil {
		return
	}

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}
	if len(data["beatmapsets"].([]interface{})) < 1 {
		*lu = ""
		*id = ""
		return
	}

	if err = updateSearchBeatmaps(data); err != nil {
		return
	}
	c := data["cursor"].(map[string]interface{})
	*lu = c["last_update"].(string)
	*id = c["_id"].(string)
	return
}

func updateMap(SET map[string]interface{}) {
	if Settings.Config.AutoDownload70FavOver {
		favCount := int(SET["favourite_count"].(float64))
		ranked := int(SET["ranked"].(float64))
		sid := strconv.Itoa(int(SET["id"].(float64)))
		if Settings.Config.Logger.ShowFavouriteCount.ALL {
			fmt.Println(sid, favCount, "favourite count")
		}
		if favCount > 70 && (ranked == 1 || ranked == 4) {
			oszFileName := sid + ".osz"
			if Settings.Config.Logger.ShowFavouriteCount.Over70 {
				fmt.Println("[U]", sid, favCount, "favourite count")
			}
			if _, err := os.Stat(oszFileName); os.IsNotExist(err) {
				fmt.Println("[C]", sid, "file dose not exist, download start")
				dl, err := DownloadBeatmap(sid, false)
				if err != nil && dl {
					fmt.Println(sid, "favourite count is 70 over but download failed.")
				}
			}
		}
	}
	//        beatmapset_id, title, title_unicode, artist, artist_unicode, creator, submitted_date,
	//        ranked, ranked_date, last_updated, play_count, bpm, tags, genre_id,
	//        genre_name, language_id, language_name, favourite_count

	Upsert(UpsertMapsSet2, []interface{}{
		SET["id"], SET["title"], SET["title_unicode"], SET["artist"], SET["artist_unicode"], SET["creator"], ToDateTime(SET["submitted_date"]),
		SET["ranked"], ToDateTime(SET["ranked_date"]), ToDateTime(SET["last_updated"]), SET["play_count"], SET["bpm"], SET["tags"], SET["favourite_count"],
	})
	for _, jz := range SET["beatmaps"].([]interface{}) {
		MAP := jz.(map[string]interface{})
		go func() {
			//		id, set_id, set_ranked, set_ranked_txt, ranked, ranked_txt, mode,
			//		mode_txt, title, title_unicode, artist, artist_unicode, version, creator, creator_id, set_submitted_date, set_last_updated,
			//		set_ranked_date, last_updated, favourite_count, set_playcount, difficulty_rating, set_bpm, bpm, ar, cs, hp,
			//		od, max_combo, playcount, passcount, total_length, hit_length, count_circles, count_spinners, count_sliders, has_storyboard,
			//		has_video, language_id, language_name, genre_id, genre_name, tags, beatmaps_count
			Upsert(UpsertMaps2, []interface{}{
				MAP["id"], MAP["beatmapset_id"], SET["ranked"], SET["status"], MAP["ranked"], MAP["status"], MAP["mode_int"],
				MAP["mode"], SET["title"], SET["title_unicode"], SET["artist"], SET["artist_unicode"], MAP["version"], SET["creator"], SET["user_id"], ToDateTime(SET["submitted_date"]), ToDateTime(SET["last_updated"]),
				ToDateTime(SET["ranked_date"]), ToDateTime(MAP["last_updated"]), SET["favourite_count"], SET["play_count"], MAP["difficulty_rating"], SET["bpm"], MAP["bpm"], MAP["ar"], MAP["cs"], MAP["drain"],
				MAP["accuracy"], MAP["max_combo"], MAP["playcount"], MAP["passcount"], MAP["total_length"], MAP["hit_length"], MAP["count_circles"], MAP["count_spinners"], MAP["count_sliders"], SET["storyboard"],
				SET["video"], SET["tags"], len(SET["beatmaps"].([]interface{})),
			})
		}()
	}
}

func updateSearchBeatmaps(data map[string]interface{}) (err error) {
	for _, v := range data["beatmapsets"].([]interface{}) {
		SET := v.(map[string]interface{})
		if Settings.Config.AutoDownload70FavOver {
			favCount := int(SET["favourite_count"].(float64))
			ranked := int(SET["ranked"].(float64))
			sid := strconv.Itoa(int(SET["id"].(float64)))
			if Settings.Config.Logger.ShowFavouriteCount.ALL {
				fmt.Println(sid, favCount, "favourite count")
			}
			if favCount > 70 && (ranked == 1 || ranked == 4) {
				oszFileName := sid + ".osz"
				if Settings.Config.Logger.ShowFavouriteCount.Over70 {
					fmt.Println("[U]", sid, favCount, "favourite count")
				}
				if _, err := os.Stat(oszFileName); os.IsNotExist(err) {
					fmt.Println("[C]", sid, "file dose not exist, download start")
					dl, err := DownloadBeatmap(sid, false)
					if err != nil && dl {
						fmt.Println(sid, "favourite count is 70 over but download failed.")
					}
				}
			}
		}
		//        beatmapset_id, title, title_unicode, artist, artist_unicode, creator, submitted_date,
		//        ranked, ranked_date, last_updated, play_count, bpm, tags, favourite_count
		Upsert(UpsertMapsSet2, []interface{}{
			SET["id"], SET["title"], SET["title_unicode"], SET["artist"], SET["artist_unicode"], SET["creator"], ToDateTime(SET["submitted_date"]),
			SET["ranked"], ToDateTime(SET["ranked_date"]), ToDateTime(SET["last_updated"]), SET["play_count"], SET["bpm"], SET["tags"], SET["favourite_count"],
		})
		for _, jz := range SET["beatmaps"].([]interface{}) {
			MAP := jz.(map[string]interface{})
			go func() {
				//		id, set_id, set_ranked, set_ranked_txt, ranked, ranked_txt, mode,
				//		mode_txt, title, title_unicode, artist, artist_unicode, version, creator, creator_id, set_submitted_date, set_last_updated,
				//		set_ranked_date, last_updated, favourite_count, set_playcount, difficulty_rating, set_bpm, bpm, ar, cs, hp,
				//		od, max_combo, playcount, passcount, total_length, hit_length, count_circles, count_spinners, count_sliders, has_storyboard,
				//		has_video, tags, beatmaps_count
				Upsert(UpsertMaps2, []interface{}{
					MAP["id"], MAP["beatmapset_id"], SET["ranked"], SET["status"], MAP["ranked"], MAP["status"], MAP["mode_int"],
					MAP["mode"], SET["title"], SET["title_unicode"], SET["artist"], SET["artist_unicode"], MAP["version"], SET["creator"], SET["user_id"], ToDateTime(SET["submitted_date"]), ToDateTime(SET["last_updated"]),
					ToDateTime(SET["ranked_date"]), ToDateTime(MAP["last_updated"]), SET["favourite_count"], SET["play_count"], MAP["difficulty_rating"], SET["bpm"], MAP["bpm"], MAP["ar"], MAP["cs"], MAP["drain"],
					MAP["accuracy"], MAP["max_combo"], MAP["playcount"], MAP["passcount"], MAP["total_length"], MAP["hit_length"], MAP["count_circles"], MAP["count_spinners"], MAP["count_sliders"], SET["storyboard"],
					SET["video"], SET["tags"], len(SET["beatmaps"].([]interface{})),
				})
			}()
		}
	}
	return
}
