package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

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

func ManualUpdateBeatmapSet(id int) {
	url := fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d", id)
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
	ms := string(body)
	if ms == "" || ms == "{\"error\":null}" || res.StatusCode != 200 {
		return
	}

	var v osu.BeatmapSets
	if err := json.Unmarshal([]byte(ms), &v); err != nil {
		fmt.Print(id, "error", err.Error())
		return
	}
	updateMap(&v)
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

func updateMap(s *osu.BeatmapSets) {
	if Settings.Config.AutoDownload70FavOver {
		favCount := s.FavouriteCount
		ranked := s.Ranked
		sid := strconv.Itoa(s.Id)
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
		s.Id, s.Title, s.TitleUnicode, s.Artist, s.ArtistUnicode, s.Creator, s.SubmittedDate,
		s.Ranked, s.RankedDate, s.LastUpdated, s.PlayCount, s.Bpm, s.Tags, s.FavouriteCount,
	})
	ch := make(chan struct{}, len(*s.Beatmaps))
	for _, m := range *s.Beatmaps {
		m := m
		go func() {
			Upsert(UpsertMaps2, []interface{}{
				m.Id, m.BeatmapsetId, s.Ranked, s.Status, m.Ranked, m.Status, m.ModeInt,
				m.Mode, s.Title, s.TitleUnicode, s.Artist, s.ArtistUnicode, m.Version, s.Creator, s.UserId, s.SubmittedDate, s.LastUpdated,
				s.RankedDate, m.LastUpdated, s.FavouriteCount, s.PlayCount, m.DifficultyRating, s.Bpm, m.Bpm, m.Ar, m.Cs, m.Drain,
				m.Accuracy, m.MaxCombo, m.Playcount, m.Passcount, m.TotalLength, m.HitLength, m.CountCircles, m.CountSpinners, m.CountSliders, s.Storyboard,
				s.Video, s.Tags, len(*s.Beatmaps),
			})
			ch <- struct{}{}
		}()
	}
	for i := 0; i < len(ch); i++ {
		<-ch
	}
}

func updateSearchBeatmaps(data *[]osu.BeatmapSets) (err error) {
	for _, s := range *data {
		if Settings.Config.AutoDownload70FavOver {
			favCount := int(s.FavouriteCount)
			ranked := int(s.Ranked)
			sid := strconv.Itoa(s.Id)
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
			s.Id, s.Title, s.TitleUnicode, s.Artist, s.ArtistUnicode, s.Creator, s.SubmittedDate,
			s.Ranked, s.RankedDate, s.LastUpdated, s.PlayCount, s.Bpm, s.Tags, s.FavouriteCount,
		})

		ch := make(chan struct{}, len(*s.Beatmaps))
		for _, m := range *s.Beatmaps {
			m := m
			go func() {
				Upsert(UpsertMaps2, []interface{}{
					m.Id, m.BeatmapsetId, s.Ranked, s.Status, m.Ranked, m.Status, m.ModeInt,
					m.Mode, s.Title, s.TitleUnicode, s.Artist, s.ArtistUnicode, m.Version, s.Creator, s.UserId, s.SubmittedDate, s.LastUpdated,
					s.RankedDate, m.LastUpdated, s.FavouriteCount, s.PlayCount, m.DifficultyRating, s.Bpm, m.Bpm, m.Ar, m.Cs, m.Drain,
					m.Accuracy, m.MaxCombo, m.Playcount, m.Passcount, m.TotalLength, m.HitLength, m.CountCircles, m.CountSpinners, m.CountSliders, s.Storyboard,
					s.Video, s.Tags, len(*s.Beatmaps),
				})
				ch <- struct{}{}
			}()
		}
		for i := 0; i < len(ch); i++ {
			<-ch
		}
	}
	return
}
