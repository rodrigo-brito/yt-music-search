package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Payload struct {
	Query   string  `json:"query"`
	Params  string  `json:"params"`
	Context Context `json:"context"`
}

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0"
	publicKey = "AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30" // public key from youtube music (it is not a leak)
	api       = "https://music.youtube.com/youtubei/v1/search"
)

type Client struct {
	Clientname    string `json:"clientName"`
	Clientversion string `json:"clientVersion"`
	Hl            string `json:"hl"`
	Gl            string `json:"gl"`
}

type Context struct {
	Client Client `json:"client"`
}

func Search(query string) {

	data, err := json.Marshal(Payload{
		Query:  query,
		Params: "EgWKAQIIAWoMEAMQBBAOEAoQBRAJ", // filter contruction https://github.com/sigma67/ytmusicapi/blob/267db615a3fcda870b36bcd83a98801e204722c9/ytmusicapi/mixins/browsing.py#L123-L164
		Context: Context{
			Client: Client{
				Hl:            "pt",
				Gl:            "BR",
				Clientname:    "WEB_REMIX",
				Clientversion: "0.1",
			},
		},
	})

	req, err := http.NewRequest("POST", fmt.Sprintf("%s?alt=json&key=%s", api, publicKey), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("origin", "https://music.youtube.com")
	req.Header.Set("referer", "https://music.youtube.com")
	req.Header.Set("content-encoding", "gzip")
	req.Header.Set("accept", "*/*")
	req.Header.Set("x-goog-authuser", "0")
	req.Header.Set("user-agent", userAgent)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func main() {
	Search("Greta Van Fleet") // check out.json
}
