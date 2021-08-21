package main

import (
	"bytes"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nerina1241/osu-beatmap-mirror-api/LoadBalancer"
	"github.com/nerina1241/osu-beatmap-mirror-api/Logger"
	"github.com/nerina1241/osu-beatmap-mirror-api/Route"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
	"github.com/nerina1241/osu-beatmap-mirror-api/middleWareFunc"
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
)

var LogIO = bytes.Buffer{}
var Logbuf = bytes.Buffer{}

func init() {
	asciiArt()
	ch := make(chan struct{})
	Settings.LoadSetting()
	go src.StartIndex()
	go src.LoadBancho(ch)
	src.ConnectMaria()
	go Logger.LoadLogger(&LogIO)
	_ = <-ch
	go src.RunGetBeatmapDataASBancho()
}

func asciiArt() {
	fmt.Println("    _   __          _                            ___    ____  ____")
	fmt.Println("   / | / /__  _____(_)___  __  ______ _____     /   |  / __ \\/  _/")
	fmt.Println("  /  |/ / _ \\/ ___/ / __ \\/ / / / __ `/ __ \\   / /| | / /_/ // /  ")
	fmt.Println(" / /|  /  __/ /  / / / / / /_/ / /_/ / / / /  / ___ |/ ____// /   ")
	fmt.Println("/_/ |_/\\___/_/  /_/_/ /_/\\__, /\\__,_/_/ /_/  /_/  |_/_/   /___/   ")
	fmt.Println("                        /____/                                    ")
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(
		// middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{echo.GET}}),
		// middleware.CORS(),
		middleware.RateLimiterWithConfig(middleWareFunc.RateLimiterConfig),
		middleware.LoggerWithConfig(middleware.LoggerConfig{Output: &LogIO}),
		middleware.RequestID(),
	)

	e.GET("/", Route.IndexPage)
	e.GET("/u", Route.UpdateBeatmap)
	e.GET("/d", Route.BeatmapDownload)
	e.GET("/download", LoadBalancer.CheckServerType)
	e.GET("/search", Route.Search)
	e.GET("/beatmapset/:sid", Route.ApiBeatmapset)

	e.Logger.Fatal(e.Start(":" + Settings.Config.Port)) // localhost:8002

}
