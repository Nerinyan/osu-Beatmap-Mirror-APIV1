package LoadBalancer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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
	switch rqdata.Server {
	case 0:
		fmt.Println("[R]", "Request Redirect to Loadbalance Downloader", rqdata.Beatmapsetid)
		return LoadBalanceDownload(c, rqdata.Beatmapsetid)
	case 1:
		fmt.Println("[R]", "Request Redirect to Main Server", rqdata.Beatmapsetid)
		return Route.DownloadBeatmapSet(c, rqdata.Beatmapsetid)
	case 2:
		fmt.Println("[R]", "Request Redirect to thftServer", rqdata.Beatmapsetid)
		return RedirectThftgrServer(c, rqdata.Beatmapsetid)
	default:
		fmt.Println("[R]", "Request Redirect to Main Server", rqdata.Beatmapsetid)
		return Route.DownloadBeatmapSet(c, rqdata.Beatmapsetid)
	}
}

func LoadBalanceDownload(c echo.Context, mid int) (err error) {
	fmt.Println("[R]", "loadbalancer Count:", Global.LoadBalance)
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
