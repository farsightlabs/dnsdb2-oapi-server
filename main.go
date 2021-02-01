// Copyright (c) 2021 by Farsight Security, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/farsightlabs/dnsdb2-oapi-proxy/pkg/proxy"
)

func main() {
	var tls bool
	var addr, remote, certFile, keyFile string
	flag.BoolVar(&tls, `tls`, false, `Use TLS`)
	flag.StringVar(&certFile, `certFile`, `cert.pem`, `TLS cert file`)
	flag.StringVar(&keyFile, `keyFile`, `cert.key`, `TLS key file`)
	flag.StringVar(&proxy.Prefix, `prefix`, `/`, `Absolute URL prefix`)
	flag.StringVar(&addr, `server`, `:8086`, `Listener address`)
	flag.StringVar(&remote, `remote`, `https://api.dnsdb.info`, `API service`)
	flag.Parse()

	var err error
	proxy.Remote, err = url.Parse(remote)
	if err != nil {
		log.Fatalf("%s: %s", remote, err)
	}

	log.Printf("Server started")

	router := proxy.NewRouter()
	server := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(router, &http2.Server{}),
	}

	if tls {
		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}
