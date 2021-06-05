package src

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
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
	fmt.Println(time.Now().UTC(), "indexing START========")

	files, err := ioutil.ReadDir(Setting.TargetDir)
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
	fmt.Println(time.Now().UTC(), "indexing END", len(FileList))
}
