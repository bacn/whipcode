//
//  Copyright 2024 whipcode.app (AnnikaV9)
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on
//  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific
//  language governing permissions and limitations under the License.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"whipcode/config"
	"whipcode/control"
	"whipcode/podman"
	"whipcode/routes"
	"whipcode/server"
)

const VERSION = "1.1.0"

func main() {
	var version, enableTLS, enableCache, enablePing, standalone bool
	var port, maxBytesSize, rlBurst, rlRefill, timeout int
	var keyFile, proxy string

	fileConfig := config.LoadConfig("config.toml")

	flag.Usage = func() {
		fmt.Printf("usage: %s [options]\n", os.Args[0])
		fmt.Println(`options:
    -h, --help              print this help message
    -v, --version           print version information
    -p, --port     PORT     port to listen on (default: 8000)
    -m, --max      BYTES    max bytes to accept (default: 1000000)
    -t, --timeout  SECONDS  timeout for execution (default: 10)
    -k, --key      FILE     master key file (default: .masterkey)
    --proxy        ADDR     reverse proxy address (default: none)
    --cache                 enable execution cache
    --tls                   enable tls
    --ping                  enable /ping endpoint
    --standalone            enable rate limiting (CHECK README)
    --burst        COUNT    rate limit burst (default: 3)
    --refill	   SECONDS  rate limit refill time (default: 1)`)
	}
	flag.BoolVar(&version, "version", false, "")
	flag.BoolVar(&version, "v", false, "")
	flag.IntVar(&port, "port", fileConfig.Port, "")
	flag.IntVar(&port, "p", fileConfig.Port, "")
	flag.IntVar(&maxBytesSize, "max", fileConfig.Max, "")
	flag.IntVar(&maxBytesSize, "m", fileConfig.Max, "")
	flag.IntVar(&timeout, "timeout", fileConfig.Timeout, "")
	flag.IntVar(&timeout, "t", fileConfig.Timeout, "")
	flag.StringVar(&keyFile, "key", fileConfig.Key, "")
	flag.StringVar(&keyFile, "k", fileConfig.Key, "")
	flag.StringVar(&proxy, "proxy", fileConfig.Proxy, "")
	flag.BoolVar(&enableCache, "cache", fileConfig.Cache, "")
	flag.BoolVar(&enableTLS, "tls", fileConfig.TLS, "")
	flag.BoolVar(&enablePing, "ping", fileConfig.Ping, "")
	flag.BoolVar(&standalone, "standalone", fileConfig.Standalone, "")
	flag.IntVar(&rlBurst, "burst", fileConfig.Burst, "")
	flag.IntVar(&rlRefill, "refill", fileConfig.Refill, "")
	flag.Parse()

	if version {
		fmt.Println(VERSION)
		return
	}

	if err := os.MkdirAll(filepath.Join(".", "run"), 0755); err != nil {
		log.Fatalf("Fatal: Could not create run directory: %v", err)
	}

	keyStore, keyAndSalt := control.InitializeKeystore(keyFile)

	scopedParams := server.ScopedMiddleWareParams{
		EnableCache:  enableCache,
		KeyAndSalt:   keyAndSalt,
		KeyStore:     keyStore,
		MaxBytesSize: maxBytesSize,
		Executor:     *podman.NewExecutor(timeout),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		server.Send(w, http.StatusNotFound, []byte(`{"detail": "not found"}`))
	})

	http.HandleFunc("POST /run", server.ScopedMiddleWare(routes.Run, scopedParams))
	http.HandleFunc("/run", func(w http.ResponseWriter, _ *http.Request) {
		server.Send(w, http.StatusMethodNotAllowed, []byte(`{"detail": "method not allowed"}`))
	})

	if enablePing {
		http.HandleFunc("/ping", routes.Ping)
	}

	rateLimiter := control.NewRateLimiter()
	if standalone {
		rateLimiter.StartCleanup()
	}

	params := server.MiddleWareParams{
		RateLimiter: rateLimiter,
		Standalone:  standalone,
		RlBurst:     rlBurst,
		RlRefill:    rlRefill,
		Proxy:       proxy,
	}

	handler := server.MiddleWare(http.DefaultServeMux, params)
	server.StartServer(port, handler, enableTLS)
}
