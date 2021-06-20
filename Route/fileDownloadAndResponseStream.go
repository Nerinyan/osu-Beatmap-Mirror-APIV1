package Route

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nerina1241/osu-beatmap-mirror-api/ConsoleLogger"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
	"github.com/pkg/errors"
)

func saveLocal(data *bytes.Buffer, path string, id int) (err error) {
	ConsoleLogger.Consolelog("Download", strconv.Itoa(id)+" | Saving beatmapsets Successful at "+Settings.Config.TargetDir+path)
	file, err := os.Create(path + ".osz.down")
	if err != nil {
		return
	}
	if file == nil {
		return errors.New("")
	}
	_, err = file.Write(data.Bytes())
	if err != nil {
		return
	}
	file.Close()

	if _, err = os.Stat(path + ".osz"); !os.IsNotExist(err) {
		err = os.Remove(path + ".osz")
		if err != nil {
			return
		}
	}
	err = os.Rename(path+".osz.down", path+".osz")
	if err != nil {
		return
	}

	src.FileList[id] = time.Now()
	ConsoleLogger.Consolelog("Download", strconv.Itoa(id)+" | Saving beatmapsets Successful at "+Settings.Config.TargetDir+path)
	return
}
func DownloadBeatmapSet(c echo.Context, mid int) (err error) {
	stringId := strconv.Itoa(mid)

	serverFileName := Settings.Config.TargetDir + "/" + stringId

	go src.ManualUpdateBeatmapSet(mid)

	rows, err := src.Maria.Query(src.GetDownloadBeatmapData, mid)
	if err != nil {
		return c.String(500, "ErrorCode: 1-1")
	}
	defer rows.Close()

	if !rows.Next() {
		return c.String(404, "please wait some second and try again or beatmap does not exist. please check beatmapset id.")
	}
	var a struct {
		Id          string
		Artist      string
		Title       string
		LastUpdated string
	}
	if err = rows.Scan(&a.Id, &a.Artist, &a.Title, &a.LastUpdated); err != nil {
		return c.String(500, "ErrorCode: 1-2")
	}

	fileName := a.Id + " " + a.Artist + " - " + a.Title + ".osz"
	ConsoleLogger.Consolelog("Check File", stringId+" | Checking if the file is latest")
	chkformat := strings.Contains(a.LastUpdated, "T")
	var lu time.Time
	if chkformat {
		lu, err = time.Parse("2006-01-02T15:04:05", a.LastUpdated)
		if err != nil {
			ConsoleLogger.WarningConsolelog("Warning", err.Error())
			return c.String(500, "ErrorCode: 1-3-1")
		}
	} else {
		lu, err = time.Parse("2006-01-02 15:04:05", a.LastUpdated)
		if err != nil {
			ConsoleLogger.WarningConsolelog("Warning", err.Error())
			return c.String(500, "ErrorCode: 1-3-2")
		}
	}

	if src.FileList[mid].Unix() >= lu.Unix() { // 맵이 최신인경우
		c.Response().Header().Set("Content-Type", "application/download")
		return c.Attachment(serverFileName+".osz", fileName)
	}

	ConsoleLogger.WarningConsolelog("Check File", stringId+" | File is not latest so redownload started")

	//==========================================
	//=        비트맵 파일이 서버에 없는경우        =
	//==========================================
	if Settings.Config.Logger.DownloadBeatmap {
		ConsoleLogger.Consolelog("Download", stringId+" | Download beatmapsets at "+Settings.Config.TargetDir)
	}
	url := "https://osu.ppy.sh/api/v2/beatmapsets/" + stringId + "/download"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return c.String(500, "ErrorCode: 2-1")
	}
	req.Header.Add("Authorization", Settings.Config.Osu.Token.TokenType+" "+Settings.Config.Osu.Token.AccessToken)

	res, err := client.Do(req)

	if err != nil {
		return c.String(500, "ErrorCode: 2-2")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		ConsoleLogger.WarningConsolelog("Download", stringId+" | Bancho returned '"+strconv.Itoa(res.StatusCode)+"' so i will retry download with another server.")
		url = "https://api.chimu.moe/v1/download/" + stringId
		client = &http.Client{}
		req, err = http.NewRequest("GET", url, nil)

		if err != nil {
			return c.String(500, "ErrorCode: 2-1-1")
		}

		res, err = client.Do(req)

		if err != nil {
			return c.String(500, "ErrorCode: 2-2-1")
		}
		defer res.Body.Close()
	}
	cLen, _ := strconv.Atoi(res.Header.Get("Content-Length"))

	c.Response().Header().Set("Content-Type", res.Header.Get("Content-Type"))
	c.Response().Header().Set("Content-Length", res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Disposition", res.Header.Get("Content-Disposition"))

	var buf = bytes.Buffer{}
	//TODO https 응답 먼저 주고 file 저장은 버퍼로 진행
	for i := 0; i < cLen; {
		var b = make([]byte, 256000)
		n, err := res.Body.Read(b)

		i += n
		buf.Write(b[:n])
		if _, err := c.Response().Write(b[:n]); err != nil {
			c.String(500, "ErrorCode: 2-4")
			return err
		}
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
	c.Response().Flush()
	return saveLocal(&buf, serverFileName, mid)
}
