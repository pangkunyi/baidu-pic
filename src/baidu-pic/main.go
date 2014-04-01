package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	API_URL = `http://image.baidu.com/i?tn=baiduimagejson&rn=2&word=%s&pn=%d`
	TPL     = `
    <html>
    <body>
    <a href="pic">
    <img src="%s" width="1024" heigth="768"/>
    </a>
    </body>
    </html>
    `
)

var (
	pageNum = 1
	port    = flag.String("port", ":8808", "server listen port")
)

type Data struct {
	LargeTnImageUrl string `json:"largeTnImageUrl"`
	ObjURL          string `json:"objURL"`
}

type Result struct {
	Datas []Data `json:"data"`
}

func search() string {
	searchUrl := fmt.Sprintf(API_URL, "tree", pageNum)
	pageNum = pageNum + 1
	if resp, err := http.Get(searchUrl); err == nil {
		defer resp.Body.Close()
		if data, err := ioutil.ReadAll(resp.Body); err == nil {
			var result Result
			if err = json.Unmarshal(data, &result); err == nil {
				return result.Datas[0].ObjURL
			}
		}
	}
	return ""
}

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	http.HandleFunc("/pic", handler)
	fmt.Println("server listen at ", *port)
	if err := http.ListenAndServe(*port, nil); err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, TPL, search())
}
