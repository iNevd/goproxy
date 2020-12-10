package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/goproxy/goproxy"
)

var addr string
var cacheDir string

func init() {
	flag.StringVar(&cacheDir, "cacheDir", "", "go modules cache dir")
	flag.StringVar(&addr, "addr", "0.0.0.0:8081", "service listen address")
	flag.Parse()

	if os.Getenv("GIT_TERMINAL_PROMPT") == "" {
		errPanic(os.Setenv("GIT_TERMINAL_PROMPT", "0"))
	}

	if os.Getenv("GIT_SSH") == "" && os.Getenv("GIT_SSH_COMMAND") == "" {
		errPanic(os.Setenv("GIT_SSH_COMMAND", "ssh -o ControlMaster=no"))
	}
}

func errPanic(err error, _ ...interface{}) {
	if err != nil {
		panic(err)
	}
}

type responseLogger struct {
	code int
	http.ResponseWriter
}

func (r *responseLogger) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

type logger struct {
	h http.Handler
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rl := &responseLogger{code: 200, ResponseWriter: w}
	l.h.ServeHTTP(rl, r)
	log.Printf("%.3fs %d %s %s\n", time.Since(start).Seconds(), rl.code, r.RemoteAddr, r.URL)
}

func main() {
	if os.Getenv("HOME") != "" {
		log.Printf("HOME %s\n", os.Getenv("HOME"))
	}

	if os.Getenv("PATH") != "" {
		log.Printf("PATH %s\n", os.Getenv("PATH"))
	}
	if os.Getenv("GOPRIVATE") != "" {
		log.Printf("GOPRIVATE %s\n", os.Getenv("GOPRIVATE"))
	}
	if os.Getenv("GOPROXY") != "" {
		log.Printf("GOPROXY %s\n", os.Getenv("GOPROXY"))
	}
	if os.Getenv("GIT_SSH_COMMAND") != "" {
		log.Printf("GIT_SSH_COMMAND %s\n", os.Getenv("GIT_SSH_COMMAND"))
	}
	log.Printf("addr %s\n", addr)

	proxy := goproxy.New()
	if cacheDir != "" {
		proxy.Cacher = &cacher.Disk{Root: cacheDir}
	}
	log.Fatal(http.ListenAndServe(addr, &logger{proxy}))
}
