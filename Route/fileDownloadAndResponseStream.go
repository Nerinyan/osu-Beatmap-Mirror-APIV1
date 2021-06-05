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
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
	"github.com/pkg/errors"
)

type Rqdata struct {
	Server       int `json:"server"`
	Beatmapsetid int `json:"beatmapsetid"`
}

func DownloadBeatmapSet(c echo.Context) (err error) {
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
	fmt.Println(rqdata)

	mid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.NoContent(500)
		return
	}

	serverFileName := src.Setting.TargetDir + "/" + c.Param("id")

	go src.ManualUpdateBeatmapSet(c.Param("id"))

	rows, err := src.Maria.Query(src.GetDownloadBeatmapData, mid)
	if err != nil {
		c.NoContent(500)
		return
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
		c.NoContent(500)
		return
	}

	fileName := a.Id + " " + a.Artist + " - " + a.Title + ".osz"
	chkformat := strings.Contains(a.LastUpdated, "T")
	if chkformat {
		lu, err := time.Parse("2006-01-02T15:04:05", a.LastUpdated)
		if err != nil {
			fmt.Println(err)
			return c.NoContent(500)
		}
		if src.FileList[mid].Unix() >= lu.Unix() { // 맵이 최신인경우
			return c.Attachment(serverFileName+".osz", fileName)
		}
	} else {
		lu, err := time.Parse("2006-01-02 15:04:05", a.LastUpdated)
		if err != nil {
			fmt.Println(err)
			return c.NoContent(500)
		}
		if src.FileList[mid].Unix() >= lu.Unix() { // 맵이 최신인경우
			return c.Attachment(serverFileName+".osz", fileName)
		}
	}

	//==========================================
	//=        비트맵 파일이 서버에 없는경우        =
	//==========================================
	fmt.Println(c.Param("id") + "file does not exist on the server, download start")
	url := "https://osu.ppy.sh/api/v2/beatmapsets/" + c.Param("id") + "/download"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		c.NoContent(500)
		return
	}
	req.Header.Add("Authorization", src.Setting.Osu.Token.TokenType+" "+src.Setting.Osu.Token.AccessToken)

	res, err := client.Do(req)

	if err != nil {
		c.NoContent(500)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		c.NoContent(404)
		return
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
			c.NoContent(500)
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

func saveLocal(data *bytes.Buffer, path string, id int) error {
	fmt.Println("beatmapSet Downloading at", src.Setting.TargetDir, path)
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
