package src

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nerina1241/osu-beatmap-mirror-api/ConsoleLogger"
	"github.com/nerina1241/osu-beatmap-mirror-api/Global"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
)

type FileIndex map[int]time.Time

var FileList = make(FileIndex)

func StartIndex() {
	for {
		FileListUpdate()
		time.Sleep(time.Second * 60 * 5)
	}
}

func FileListUpdate() {
	ConsoleLogger.Consolelog("Indexing", "File indexing has been started")

	files, err := ioutil.ReadDir(Settings.Config.TargetDir)
	if err != nil {
		panic(err)
	}

	tmp := make(FileIndex)
	for _, file := range files {
		if sid, err := strconv.Atoi(strings.Replace(file.Name(), ".osz", "", -1)); err == nil {
			tmp[sid] = file.ModTime()
		}
	}
	FileList = tmp
	Global.IndexCount = len(FileList)
	sTotalIndex := strconv.Itoa(Global.IndexCount)
	indexTotalSize, err := DirSize(Settings.Config.TargetDir)
	Global.IndexSize = indexTotalSize / 1024 / 1024 / 1024
	ConsoleLogger.Consolelog("Indexing", "File indexing done! "+sTotalIndex+" files("+strconv.FormatInt(Global.IndexSize, 10)+"GB) are indexed.")
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
