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

package proxy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	dnsdbv2 "github.com/dnsdb/go-dnsdb/pkg/dnsdb/v2"
	"github.com/dnsdb/go-dnsdb/pkg/dnsdb/v2/saf"
)

var (
	Prefix = "/"
	Remote = dnsdbv2.DefaultDnsdbServer
	Client = http.DefaultClient
)

func copyHeader(from, to http.Header) {
	for k, v := range from {
		to[k] = v
	}
}

func normalizePrefix(prefix string) string {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	if strings.HasSuffix(prefix, "/") {
		prefix = prefix[:len(prefix)-1]
	}

	return prefix
}

func proxyURL(base, proxied *url.URL) *url.URL {
	res := new(url.URL)
	*res = *base
	res.Path = proxied.Path
	prefix := normalizePrefix(Prefix)
	if strings.HasPrefix(res.Path, prefix) {
		res.Path = res.Path[len(prefix):]
	}
	res.RawQuery = proxied.Query().Encode()
	return res
}

func proxyRequest(req *http.Request) (*http.Request, error) {
	proxyReq, err := http.NewRequest(http.MethodGet, proxyURL(Remote, req.URL).String(), req.Body)
	if err != nil {
		return nil, err
	}
	copyHeader(req.Header, proxyReq.Header)
	proxyReq.Header.Set("Accept", dnsdbv2.ContentType)
	proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)
	return proxyReq.WithContext(req.Context()), nil
}

func ProxySafHandler(w http.ResponseWriter, req *http.Request) {
	proxyReq, err := proxyRequest(req)
	if err != nil {
		ErrorHandler(err, w)
		return
	}

	resp, err := Client.Do(proxyReq)
	if err != nil {
		ErrorHandler(err, w)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ErrorHandler(err, w)
			return
		}

		copyHeader(resp.Header, w.Header())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(b)
		return
	}

	copyHeader(resp.Header, w.Header())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Trailer", "Success")

	res := []json.RawMessage{}
	s := saf.Stream{}
	s.Run(req.Context(), resp.Body)
	defer s.Close()
	for rawMsg := range s.Ch() {
		res = append(res, rawMsg)
	}
	if s.Err() != nil {
		if errors.Is(s.Err(), saf.ErrStreamLimited) {
			w.Header().Set("Limited", "True")
		} else {
			ErrorHandler(s.Err(), w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	if err = e.Encode(res); err != nil {
		// This will never happen. The json.RawMessages have already been validated.
		panic(err)
	}
	w.Header().Set("Success", "True")
}

func ProxyPassHandler(w http.ResponseWriter, req *http.Request) {
	proxyReq, err := proxyRequest(req)
	if err != nil {
		ErrorHandler(err, w)
		return
	}

	resp, err := Client.Do(proxyReq)
	if err != nil {
		ErrorHandler(err, w)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorHandler(err, w)
		return
	}

	copyHeader(resp.Header, w.Header())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Trailer", "Success")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
	w.Header().Set("Success", "True")
}
