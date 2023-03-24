package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	listen    string
	outputDir string
	verbose   bool
)

func main() {
	flag.BoolVar(&verbose, "v", false, "print request to stdout")
	flag.StringVar(&listen, "listen", ":8080", "listen address and port")
	flag.StringVar(&outputDir, "output-dir", "", "directory to store the request")
	flag.Parse()

	if outputDir != "" {
		os.MkdirAll(outputDir, 0755)
	}

	server := &http.Server{}
	server.Addr = listen
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ww io.Writer = os.Stdout
		if outputDir != "" {
			ts := strconv.Itoa(int(time.Now().UnixNano()))
			f, err := os.Create(outputDir + "/" + ts + ".http")
			if err != nil {
				log.Print(err)
			}
			defer f.Close()
			ww = f
		}
		r.Header.Write(ww)
		ww.Write([]byte("\r\n"))
		_, _ = io.Copy(ww, r.Body)
		ww.Write([]byte("\r\n"))
		w.WriteHeader(http.StatusOK)
	})

	log.Print("listening", listen)
	_ = server.ListenAndServe()
}
