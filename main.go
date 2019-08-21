package main

import (
	"flag"
	"github.com/goproxy/goproxy"
	"github.com/goproxy/goproxy/cacher"
	"log"
	"net/http"
	"os"
	"time"
)

var listen string
var cacheDir string
var proxyHost string
var excludeHost string

func init() {
	flag.StringVar(&excludeHost, "exclude", "", "exclude host pattern")
	flag.StringVar(&proxyHost, "proxy", "", "next hop proxy for go modules")
	flag.StringVar(&cacheDir, "cacheDir", "", "go modules cache dir")
	flag.StringVar(&listen, "listen", "0.0.0.0:8081", "service listen address")
	flag.Parse()

	if os.Getenv("GIT_TERMINAL_PROMPT") == "" {
		os.Setenv("GIT_TERMINAL_PROMPT", "0")
	}

	if os.Getenv("GIT_SSH") == "" && os.Getenv("GIT_SSH_COMMAND") == "" {
		os.Setenv("GIT_SSH_COMMAND", "ssh -o ControlMaster=no")
	}

	if excludeHost != "" {
		os.Setenv("GOPRIVATE", excludeHost)
	}
	excludeHost = os.Getenv("GOPRIVATE")

	if proxyHost != "" {
		os.Setenv("GOPROXY", proxyHost)
	}
	proxyHost = os.Getenv("GOPROXY")
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
	if os.Getenv("GOPRIVATE") != "" {
		log.Printf("ExcludeHost %s\n", os.Getenv("GOPRIVATE"))
	}
	if os.Getenv("GOPROXY") != "" {
		log.Printf("ProxyHost %s\n", os.Getenv("GOPROXY"))
	}
	if os.Getenv("GIT_SSH_COMMAND") != "" {
		log.Printf("GIT_SSH_COMMAND %s\n", os.Getenv("GIT_SSH_COMMAND"))
	}
	log.Printf("Listen %s\n", listen)

	proxy := goproxy.New()
	proxy.Cacher = &cacher.Disk{Root: cacheDir}
	log.Fatal(http.ListenAndServe(listen, &logger{proxy}))
}
