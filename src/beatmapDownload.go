package src

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
)

func DownloadBeatmap(id string, whitName bool) (b bool, err error) {
	url := "https://osu.ppy.sh/api/v2/beatmapsets/" + id + "/download"
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
		defer res.Body.Close()

		return
	}
	defer res.Body.Close()

	filename := res.Header.Get("Content-Disposition")
	filename = strings.TrimLeft(filename, "attachment;filename=\"")
	filename = strings.TrimRight(filename, "\";")

	if res.StatusCode != 200 {
		return
	}
	_, err = os.Stat(Settings.Config.TargetDir)
	if os.IsNotExist(err) {
		fmt.Println("Folder does not exist. i will make new.")
		err = os.Mkdir(Settings.Config.TargetDir, 0755)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	if Settings.Config.Logger.DownloadBeatmap {
		fmt.Println("[D]", id, "beatmapSets Downloaded at "+Settings.Config.TargetDir)
	}

	// Create the file
	if !whitName {
		filename = id + ".osz"
	}
	out, err := os.Create(Settings.Config.TargetDir + "/" + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	return

}
