package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nerina1241/osu-beatmap-mirror-api/Route"
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
)

var LogIO = bytes.Buffer{}

func init() {
	ch := make(chan struct{})

	src.LoadSetting()
	go src.LoadBancho(ch)
	src.ConnectMaria()

}

func main() {
	src.RunGetBeatmapDataASBancho()
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/", func(c echo.Context) error {
		i := c.QueryParam("s")

		if src.ManualUpdateBeatmapSet(i) != nil {
			return c.JSON(404, `{"success":false,"message":"bancho return null or server error"}`)
		}
		fmt.Print(" Alive - ", i)
		return c.JSON(200, `{"success":true}`)
	})
	e.GET("/d", func(c echo.Context) error {
		whitName := true
		setid := c.QueryParam("s")
		if setid == "" {
			return c.HTML(400, `{"success":false,"message":"parm 's=int' is null <br> 'name=bool' 123456.osz"}`)
		}
		if _, err := strconv.Atoi(setid); err != nil {
			return c.HTML(400, `{"success":false,"message":"parm data is not int"}`)
		}
		if c.QueryParam("name") == "false" {
			whitName = false
		}
		if !downloadBeatmap(setid, whitName) {
			return c.HTML(400, `{"success":false,"message":"fail to download"}`)
		}

		return c.HTML(200, `{"success":true}`)
	})
	e.Router().Add("GET", "/download/:id", Route.DownloadBeatmapSet)

	fmt.Println("Ready API Server")

	e.Logger.Fatal(e.Start(":" + src.Setting.Port)) // localhost:8002

}
func downloadBeatmap(id string, whitName bool) bool {
	url := "https://osu.ppy.sh/api/v2/beatmapsets/" + id + "/download"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Add("Authorization", src.Setting.Osu.Token.TokenType+" "+src.Setting.Osu.Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		defer res.Body.Close()

		return false
	}
	defer res.Body.Close()

	filename := res.Header.Get("Content-Disposition")
	filename = strings.TrimLeft(filename, "attachment;filename=\"")
	filename = strings.TrimRight(filename, "\";")

	if res.StatusCode != 200 {
		return false
	}
	_, err = os.Stat(src.Setting.TargetDir)
	if os.IsNotExist(err) {
		fmt.Println("Folder does not exist.")
		err = os.Mkdir(src.Setting.TargetDir, 0755)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}
	fmt.Println("beatmapSets Downloaded at " + src.Setting.TargetDir)

	// Create the file
	if !whitName {
		filename = id + ".osz"
	}
	out, err := os.Create(src.Setting.TargetDir + "/" + filename)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true

}
