package Route

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
)

func BeatmapDownload(c echo.Context) error {
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
}
