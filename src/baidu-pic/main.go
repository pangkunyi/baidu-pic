package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

const (
	API_URL = `http://image.soso.com/pics?query=%s&mood=0&picformat=0&mode=1&di=2&w=05009900&dr=1&_asf=image.soso.com&_ast=1397814386&dm=11&leftp=44230501&cwidth=1024&cheight=768&start=%d&reqType=ajax&tn=0`
	TPL     = `
    <html>
    <body>
    <a href="/pic">
    <img src="/img?url=%s" width="1024" heigth="768"/>
    </a>
    <div>%s</div>
    </body>
    </html>
    `
)

var (
	pageNum  = 1
	keywords = make(map[string]int)
	keyword  = "树"
	port     = flag.String("port", ":8808", "server listen port")
	client   http.Client
)

type Data struct {
	LargeTnImageUrl string `json:"sohu_image"`
}

type Result struct {
	Datas []Data `json:"items"`
}

func links() string {
	links := ""
	for k, v := range keywords {
		links = links + `<a href="/pic/` + k + `">` + k + `-` + strconv.Itoa(v) + `</a>...`
	}
	return links
}

func search() string {
	searchUrl := fmt.Sprintf(API_URL, url.QueryEscape(keyword), keywords[keyword])
	keywords[keyword] = keywords[keyword] + 1
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		return ""
	}
	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()
		if data, err := ioutil.ReadAll(resp.Body); err == nil {
			var result Result
			if err = json.Unmarshal(data, &result); err == nil {
				return result.Datas[0].LargeTnImageUrl
			}
		}
	}
	return ""
}

func main() {
	keywords[keyword] = 0
	keywords["猫"] = 0
	keywords["狗"] = 0
	keywords["蜗牛"] = 0
	keywords["卡车"] = 0
	keywords["汽车"] = 0
	keywords["大象"] = 0
	keywords["猴子"] = 0
	keywords["菠萝"] = 0
	keywords["字母"] = 0
	keywords["数字"] = 0
	flag.Parse()
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client = http.Client{Jar: jar}
	http.HandleFunc("/pic", handler)
	http.HandleFunc("/pic/", handler)
	http.HandleFunc("/img", img)
	fmt.Println("server listen at ", *port)
	if err := http.ListenAndServe(*port, nil); err != nil {
		panic(err)
	}
}

func img(w http.ResponseWriter, r *http.Request) {
	u := r.FormValue("url")
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Fprintf(w, "err:%#v", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "err:%#v", err)
		return
	}
	defer resp.Body.Close()
	w.Header().Add("Content-Type", "image/jpeg")
	io.Copy(w, resp.Body)
}
func handler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 5 {
		keyword = r.URL.Path[5:]
		if _, ok := keywords[keyword]; !ok {
			keywords[keyword] = 0
		}
	}
	u := search()
	if len(u) > 200 {
		fmt.Fprintf(w, u)
		return
	}
	fmt.Fprintf(w, TPL, url.QueryEscape(u), links())
}
