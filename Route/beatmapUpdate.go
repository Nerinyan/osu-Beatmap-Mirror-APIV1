package Route

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
)

func UpdateBeatmap(c echo.Context) error {
	k := c.QueryParam("k")
	if k != Settings.Config.Key {
		return c.String(404, "ErrorCode: -1")
	}
	i := c.QueryParam("s")
	ii, error := strconv.Atoi(c.QueryParam("s"))
	if error != nil {
		return c.String(404, "ErrorCode: -2")
	}
	if src.ManualUpdateBeatmapSet(ii) != nil {
		return c.JSON(404, `{"success":false,"message":"bancho return null or server error"}`)
	}
	fmt.Println(" Alive - ", i)
	return c.JSON(200, `{"success":true}`)
}
