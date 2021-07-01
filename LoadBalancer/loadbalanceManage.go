package LoadBalancer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nerina1241/osu-beatmap-mirror-api/ConsoleLogger"
	"github.com/nerina1241/osu-beatmap-mirror-api/Global"
	"github.com/nerina1241/osu-beatmap-mirror-api/Route"
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
	fmt.Println("server req", rqdata.Server, "bsid", rqdata.Beatmapsetid)
	switch rqdata.Server {
	case 0:
		ConsoleLogger.LoadBConsolelog("LoadBalance", strconv.Itoa(rqdata.Beatmapsetid)+" | Request Redirect to Loadbalance Downloader")
		return LoadBalanceDownload(c, rqdata.Beatmapsetid)
	case 1:
		ConsoleLogger.LoadBConsolelog("LoadBalance", strconv.Itoa(rqdata.Beatmapsetid)+" | Request Redirect to Main Server")
		return Route.DownloadBeatmapSet(c, rqdata.Beatmapsetid)
	case 2:
		ConsoleLogger.LoadBConsolelog("LoadBalance", strconv.Itoa(rqdata.Beatmapsetid)+" | Request Redirect to thftServer")
		return RedirectThftgrServer(c, rqdata.Beatmapsetid)
	default:
		ConsoleLogger.LoadBConsolelog("LoadBalance", strconv.Itoa(rqdata.Beatmapsetid)+" | Request Redirect to Main Server")
		return Route.DownloadBeatmapSet(c, rqdata.Beatmapsetid)
	}
}

func LoadBalanceDownload(c echo.Context, mid int) (err error) {
	ConsoleLogger.LoadBConsolelog("LoadBalance", "Loadbalance Count "+strconv.Itoa(Global.LoadBalance))
	switch Global.LoadBalance {
	case 0:
		Global.LoadBalance = Global.LoadBalance + 1
		return Route.DownloadBeatmapSet(c, mid)
	case 1:
		Global.LoadBalance = 0
		return RedirectThftgrServer(c, mid)
	default:
		return c.String(505, "ErrorCode: 0-1")
	}
}

func RedirectThftgrServer(c echo.Context, mid int) error {
	bid := strconv.Itoa(mid)
	URL := "https://xiiov.com/d/" + bid
	return c.Redirect(http.StatusPermanentRedirect, string(URL))
}
