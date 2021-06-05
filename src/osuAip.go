package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
)

func LoadBancho(ch chan struct{}) {
	b := false
	checkUpdatable := Settings.Config.Osu.Token.UpdatedAt + Settings.Config.Osu.Token.ExpiresIn - time.Now().Unix()
	if checkUpdatable > 3600 {
		fmt.Println("bancho - token Alive")
		ch <- struct{}{}
		time.Sleep(time.Second * time.Duration(checkUpdatable-100))
	}

	for {
		err := login()
		if err != nil {
			fmt.Println("fail Get bancho Token")
			panic(err)
		}
		if !b {
			b = true
			ch <- struct{}{}
		}
		fmt.Println("successful Get bancho Token")
		Settings.Config.Osu.Token.UpdatedAt = time.Now().Unix()
		Settings.Config.Save()
		time.Sleep(time.Second * 60 * 60 * 20)
	}

}

func login() error {
	url := "https://osu.ppy.sh/oauth/token"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("username", Settings.Config.Osu.Username)
	_ = writer.WriteField("password", Settings.Config.Osu.Passwd)
	_ = writer.WriteField("grant_type", "password")
	_ = writer.WriteField("client_id", "5")
	_ = writer.WriteField("client_secret", "FGc9GAtyHzeQDshWP5Ah7dega8hJACAJpQtw6OXk")
	_ = writer.WriteField("scope", "*")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {

		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return json.Unmarshal(body, &Settings.Config.Osu.Token)
}
