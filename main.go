package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nerina1241/osu-beatmap-mirror-api/Route"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
)

var LogIO = bytes.Buffer{}

func init() {
	ch := make(chan struct{})

	Settings.LoadSetting()
	go src.LoadBancho(ch)
	src.ConnectMaria()

}

func main() {
	src.RunGetBeatmapDataASBancho()
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowMethods: []string{echo.GET}}))

	e.GET("/u", func(c echo.Context) error {
		k := c.QueryParam("k")
		if k != Settings.Config.Key {
			return c.String(404, "ErrorCode: -1")
		}

		i := c.QueryParam("s")
		if src.ManualUpdateBeatmapSet(i) != nil {
			return c.JSON(404, `{"success":false,"message":"bancho return null or server error"}`)
		}
		fmt.Print(" Alive - ", i)
		return c.JSON(200, `{"success":true}`)
	})

	e.GET("/d", func(c echo.Context) error {
		k := c.QueryParam("k")
		if k != Settings.Config.Key {
			return c.String(404, "ErrorCode: -1")
		}

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

		b, err := src.DownloadBeatmap(setid, whitName)
		if err != nil && b {
			return c.HTML(400, `{"success":false,"message":"fail to download"}`)
		}

		return c.HTML(200, `{"success":true}`)
	})

	e.GET("/download", Route.CheckServerType)

	fmt.Println("Ready API Server")

	e.Logger.Fatal(e.Start(":" + Settings.Config.Port)) // localhost:8002

}
