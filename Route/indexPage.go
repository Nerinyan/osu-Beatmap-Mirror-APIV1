package Route

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nerina1241/osu-beatmap-mirror-api/ConsoleLogger"
	"github.com/nerina1241/osu-beatmap-mirror-api/Global"
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
)

type dbData struct {
	MapsInDB            int
	BiggestBeatmapsetID int
}

type outputData struct {
	MapsInDB            int
	CachedMaps          int
	CachedMapsSizeGB    int64
	CachedMapsSizeMB    int64
	BiggestBeatmapsetID int
}

func IndexPage(c echo.Context) (err error) {
	rows, err := src.Maria.Query("select count(beatmap_id) as MapsInDB, max(beatmapset_id) as BiggestBeatmapsetID from BeatmapMirror.beatmap;")

	if err != nil {
		ConsoleLogger.WarningConsolelog("Error", err.Error())
		c.NoContent(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var (
		MapsInDB            int
		BiggestBeatmapsetID int
	)
	for rows.Next() {
		var d dbData
		err = rows.Scan(&d.MapsInDB, &d.BiggestBeatmapsetID)
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return
		}
		MapsInDB = d.MapsInDB
		BiggestBeatmapsetID = d.BiggestBeatmapsetID
	}
	var output outputData
	output.MapsInDB = MapsInDB
	output.CachedMaps = Global.IndexCount
	output.CachedMapsSizeGB = Global.IndexSize
	output.CachedMapsSizeMB = Global.IndexSize * 1024
	output.BiggestBeatmapsetID = BiggestBeatmapsetID

	return c.JSON(http.StatusOK, output)

}
