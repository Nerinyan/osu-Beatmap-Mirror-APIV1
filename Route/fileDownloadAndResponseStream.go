package Route

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nerina1241/osu-beatmap-mirror-api/Global"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
	"github.com/pkg/errors"
)

type Rqdata struct {
	Server       int `json:"server"`
	Beatmapsetid int `json:"beatmapsetid"`
}

func CheckServerType(c echo.Context) (err error) {
	base64Setting := c.QueryParam("b")
	StringJsonSetting, err := base64.StdEncoding.DecodeString(base64Setting)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	JsonData := []byte(StringJsonSetting)
	var rqdata Rqdata
	err = json.Unmarshal(JsonData, &rqdata)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch rqdata.Server {
	case 0:
		return LoadBalanceDownload(c, rqdata.Beatmapsetid)
	case 1:
		return DownloadBeatmapSet(c, rqdata.Beatmapsetid)
	case 2:
		return RedirectThftgrServer(c, rqdata.Beatmapsetid)
	default:
		return DownloadBeatmapSet(c, rqdata.Beatmapsetid)
	}
}

func LoadBalanceDownload(c echo.Context, mid int) (err error) {
	switch Global.LoadBalance {
	case 0:
		Global.LoadBalance = Global.LoadBalance + 1
		return DownloadBeatmapSet(c, mid)
	case 1:
		Global.LoadBalance = 0
		return RedirectThftgrServer(c, mid)
	default:
		return c.String(505, "ErrorCode: 0-1")
	}
}

func DownloadBeatmapSet(c echo.Context, mid int) (err error) {
	stringId := strconv.Itoa(mid)

	serverFileName := Settings.Config.TargetDir + "/" + stringId

	go src.ManualUpdateBeatmapSet(stringId)

	rows, err := src.Maria.Query(src.GetDownloadBeatmapData, mid)
	if err != nil {
		return c.String(500, "ErrorCode: 1-1")
	}
	defer rows.Close()
	if !rows.Next() {
		return c.String(404, "please wait some second and try again or later")
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
	chkformat := strings.Contains(a.LastUpdated, "T")
	if chkformat {
		lu, err := time.Parse("2006-01-02T15:04:05", a.LastUpdated)
		if err != nil {
			fmt.Println(err)
			return c.String(500, "ErrorCode: 1-3-1")
		}
		if src.FileList[mid].Unix() >= lu.Unix() { // 맵이 최신인경우
			return c.Attachment(serverFileName+".osz", fileName)
		}
	} else {
		lu, err := time.Parse("2006-01-02 15:04:05", a.LastUpdated)
		if err != nil {
			fmt.Println(err)
			return c.String(500, "ErrorCode: 1-3-2")
		}
		if src.FileList[mid].Unix() >= lu.Unix() { // 맵이 최신인경우
			return c.Attachment(serverFileName+".osz", fileName)
		}
	}

	//==========================================
	//=        비트맵 파일이 서버에 없는경우        =
	//==========================================
	if Settings.Config.Logger.DownloadBeatmap {
		fmt.Println("[d] " + stringId + "file does not exist on the server, download start")
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
		fmt.Println("Bancho returned", res.StatusCode, "so i will try to use another server")
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

	c.Response().Header().Set("Content-Type", "application/zip")
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

func RedirectThftgrServer(c echo.Context, mid int) error {
	bid := strconv.Itoa(mid)
	URL := "https://xiiov.com/d/" + bid
	return c.Redirect(http.StatusPermanentRedirect, string(URL))
}

func saveLocal(data *bytes.Buffer, path string, id int) error {
	fmt.Println("beatmapSet Downloading at", Settings.Config.TargetDir, path)
	file, err := os.Create(path + ".osz.down")
	if err != nil {
		return err
	}
	if file == nil {
		return errors.New("")
	}
	_, err = file.Write(data.Bytes())
	if err != nil {
		return err
	}
	file.Close()

	err = os.Rename(path+".osz.down", path+".osz")
	if err != nil {
		return err
	}
	src.FileList[id] = time.Now()
	return nil
}
