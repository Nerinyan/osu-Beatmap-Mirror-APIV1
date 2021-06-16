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
	fmt.Println("bancho - token Dead")
	fmt.Println("bancho - start refresh bancho token")

	for {
		fmt.Println("bancho - login try")
		err := login()
		if err != nil {
			fmt.Println("fail Get bancho Token")
			panic(err)
		}
		fmt.Println("bancho - login success")
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

// func login2() error {
// 	c := &http.Client{Timeout: time.Second * 10}
// 	URL := "https://osu.ppy.sh/oauth/token"

// 	v := url.Values{}
// 	// v.Set("username", Settings.Config.Osu.Username)
// 	// v.Set("password", Settings.Config.Osu.Passwd)
// 	v.Set("grant_type", "password")
// 	v.Set("client_id", "5")
// 	v.Set("client_secret", "FGc9GAtyHzeQDshWP5Ah7dega8hJACAJpQtw6OXk")
// 	v.Set("scope", "*")

// 	req, err := http.NewRequest("POST", URL, strings.NewReader(v.Encode()))
// 	req.SetBasicAuth(Settings.Config.Osu.Username, Settings.Config.Osu.Passwd)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	fmt.Println(req)

// 	resp, err := c.Do(req)
// 	if err != nil {
// 		fmt.Println("에러", err)
// 		return err
// 	}
// 	fmt.Println(resp)

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	fmt.Println(body)

// 	return json.Unmarshal(body, &Settings.Config.Osu.Token)
// }

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
		fmt.Println(err)
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
		fmt.Println(err)
		return
	}
	return json.Unmarshal(body, &Settings.Config.Osu.Token)
}
