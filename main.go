package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/namsral/flag"
)

type (
	Sherpa struct {
		upstream *url.URL
	}
)

func NewSherpa(upstream *url.URL) *Sherpa {
	return &Sherpa{upstream}
}

func (s *Sherpa) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req *http.Request
	var err error

	if err = r.ParseForm(); err != nil {
		log.Println(err)
	}

	upstream := *s.upstream
	upstream.Path = r.RequestURI
	req, err = http.NewRequest(r.Method, upstream.String(), r.Body)
	if err != nil {
		log.Println(err)
	}
	req.Header = r.Header
	req.Header.Set("User-Agent", "sinar-sherpa")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	for k, l := range resp.Header {
		if k == "Set-Cookie" {
			continue
		}

		for _, v := range l {
			w.Header().Set(k, v)
		}
	}
	w.Header().Set("Server", "sinar-sherpa")
	w.Write(body)
}

var host, ui, upstreams string

func init() {
	flag.StringVar(&upstreams, "upstreams", "", "upstreams service url (example: https://popit.mysociety.org)")
	flag.StringVar(&host, "http", "0.0.0.0:8080", "<addr>:<port> to listen on")
	flag.StringVar(&ui, "ui", "", "path to html ui")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer func() {
		// let's skip the gory details, shall we
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(ui)))
	// Let's proxy all of these babies
	for _, u := range strings.Split(upstreams, ",") {
		u = strings.TrimSpace(u)

		uu, err := url.Parse(u)
		if err != nil {
			fmt.Printf("received invalid upstream %s\n", u)
			os.Exit(1)
		}

		var path string
		if uu.Path == "" {
			path = "/"
		} else {
			path = uu.Path
		}
		uu.Path = "/"

		p := NewSherpa(uu)
		fmt.Printf("proxying %s => \"%s\"\n", path, u)
		mux.Handle(path, p)
	}

	s := negroni.New(negroni.NewLogger(), negroni.NewRecovery())
	s.UseHandler(mux)
	s.Run(host)
}
