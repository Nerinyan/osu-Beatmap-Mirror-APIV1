package ConsoleLogger

import (
	"fmt"
	"time"
)

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func Consolelog(pType, pText string) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	now = now + " (UTC)"
	p := fmt.Sprintln(now, string(colorBlue)+"["+pType+"]", string(colorReset)+pText)
	fmt.Print(p)
}

func DangersConsolelog(pType, pText string) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	now = now + " (UTC)"
	p := fmt.Sprintln(now, string(colorRed)+"["+pType+"]", string(colorReset)+pText)
	fmt.Print(p)
}

func WarningConsolelog(pType, pText string) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	now = now + " (UTC)"
	p := fmt.Sprintln(now, string(colorYellow)+"["+pType+"]", string(colorReset)+pText)
	fmt.Print(p)
}

func GoodConsolelog(pType, pText string) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	now = now + " (UTC)"
	p := fmt.Sprintln(now, string(colorGreen)+"["+pType+"]", string(colorReset)+pText)
	fmt.Print(p)
}

func LoadBConsolelog(pType, pText string) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	now = now + " (UTC)"
	p := fmt.Sprintln(now, string(colorPurple)+"["+pType+"]", string(colorReset)+pText)
	fmt.Print(p)
}

func UpdateLConsolelog(pType, pText string) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	now = now + " (UTC)"
	p := fmt.Sprintln(now, string(colorCyan)+"["+pType+"]", string(colorReset)+pText)
	fmt.Print(p)
}
