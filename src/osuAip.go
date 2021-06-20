package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/nerina1241/osu-beatmap-mirror-api/ConsoleLogger"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
)

func LoadBancho(ch chan struct{}) {
	b := false
	checkUpdatable := Settings.Config.Osu.Token.UpdatedAt + Settings.Config.Osu.Token.ExpiresIn - time.Now().Unix()
	if checkUpdatable > 3600 {
		ConsoleLogger.Consolelog("Bancho", "Api Token Alived!")
		ch <- struct{}{}
		time.Sleep(time.Second * time.Duration(checkUpdatable-100))
	}
	ConsoleLogger.DangersConsolelog("Bancho", "Api Token Dead")
	ConsoleLogger.WarningConsolelog("Bancho", "Refreshing Bancho Api Token...")
	for {
		ConsoleLogger.GoodConsolelog("Bancho", "Api Token Generate - Login Tryed")
		err := login()
		if err != nil {
			ConsoleLogger.DangersConsolelog("Bancho", "Api Token Generate - Failed To Login")
			panic(err)
		}
		ConsoleLogger.GoodConsolelog("Bancho", "Api Token Generate - Login Successful")
		if !b {
			b = true
			ch <- struct{}{}
		}
		ConsoleLogger.Consolelog("Bancho", "Succesfully Generated Bancho Api Token")
		Settings.Config.Osu.Token.UpdatedAt = time.Now().Unix()
		Settings.Config.Save()
		time.Sleep(time.Second * 60 * 60 * 20) //20 hours
	}

}

func login() (err error) {
	url := "https://osu.ppy.sh/oauth/token"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("username", string(Settings.Config.Osu.Username))
	_ = writer.WriteField("password", string(Settings.Config.Osu.Passwd))
	_ = writer.WriteField("grant_type", "password")
	_ = writer.WriteField("client_id", "5")
	_ = writer.WriteField("client_secret", "FGc9GAtyHzeQDshWP5Ah7dega8hJACAJpQtw6OXk")
	_ = writer.WriteField("scope", "*")
	err = writer.Close()
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {

		fmt.Println("err", err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ConsoleLogger.WarningConsolelog("Warning", err.Error())
		return
	}
	return json.Unmarshal(body, &Settings.Config.Osu.Token)
}
