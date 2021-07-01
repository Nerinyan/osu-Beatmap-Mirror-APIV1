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
		ConsoleLogger.GoodConsolelog("Bancho", "Api Token Alive")
		ch <- struct{}{}
		time.Sleep(time.Second * time.Duration(checkUpdatable-100))
	}
	ConsoleLogger.DangersConsolelog("Bancho", "Api Token Dead")
	ConsoleLogger.WarningConsolelog("Bancho", "Refreshing Bancho Api Token...")
	for {
		ConsoleLogger.GoodConsolelog("Bancho", "Api Token Generate - Login Tryed")
		if err := login(true); err != nil {
			ConsoleLogger.DangersConsolelog("Bancho", "Api Token Generate - Failed To Login")
			if er := login(false); er != nil {
				panic("fail LOGIN bancho - warn")
			}
			ConsoleLogger.GoodConsolelog("Bancho", "Api Token Generate - Login Successful")
		} else {
			ConsoleLogger.Consolelog("Bancho", "Succesfully Generated Bancho Api Token")
		}
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

func login(refresh bool) (err error) {
	url := "https://osu.ppy.sh/oauth/token"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("client_id", "5")
	_ = writer.WriteField("client_secret", "FGc9GAtyHzeQDshWP5Ah7dega8hJACAJpQtw6OXk")
	_ = writer.WriteField("scope", "*")

	if refresh {
		_ = writer.WriteField("grant_type", "refresh_token")
		_ = writer.WriteField("refresh_token", Settings.Config.Osu.Token.RefreshToken)
	} else {
		_ = writer.WriteField("username", Settings.Config.Osu.Username)
		_ = writer.WriteField("password", Settings.Config.Osu.Passwd)
		_ = writer.WriteField("grant_type", "password")
	}

	err = writer.Close()
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if refresh {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", Settings.Config.Osu.Token.TokenType, Settings.Config.Osu.Token.AccessToken))
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(body, &Settings.Config.Osu.Token)
}
