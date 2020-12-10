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

var listen string
var cacheDir string
var proxyHost string
var excludeHost string
var envPath string

func init() {
	flag.StringVar(&excludeHost, "exclude", "", "exclude host pattern")
	flag.StringVar(&proxyHost, "proxy", "", "next hop proxy for go modules")
	flag.StringVar(&cacheDir, "cacheDir", "", "go modules cache dir")
	flag.StringVar(&listen, "listen", "0.0.0.0:8081", "service listen address")
	flag.StringVar(&envPath, "path", "", "PATH ENV")
	flag.Parse()

	if os.Getenv("GIT_TERMINAL_PROMPT") == "" {
		errPanic(os.Setenv("GIT_TERMINAL_PROMPT", "0"))
	}

	if os.Getenv("GIT_SSH") == "" && os.Getenv("GIT_SSH_COMMAND") == "" {
		errPanic(os.Setenv("GIT_SSH_COMMAND", "ssh -o ControlMaster=no"))
	}

	if os.Getenv("HOME") == "" {
		errPanic(os.Setenv("HOME", path.Dir(os.Args[0])))
	}

	if excludeHost != "" {
		errPanic(os.Setenv("GOPRIVATE", excludeHost))
	}
	excludeHost = os.Getenv("GOPRIVATE")

	if proxyHost != "" {
		errPanic(os.Setenv("GOPROXY", proxyHost))
	}
	proxyHost = os.Getenv("GOPROXY")

	if envPath != "" {
		errPanic(os.Setenv("PATH", envPath))
	}
	envPath = os.Getenv("PATH")
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
	log.Fatal(http.ListenAndServe(listen, &logger{proxy}))
}
