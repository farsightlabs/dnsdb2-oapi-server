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
	"context"
	"net/http"
	"net/url"
	"testing"

	v2 "github.com/dnsdb/go-dnsdb/pkg/dnsdb/v2"

	. "github.com/onsi/gomega"
)

func TestCopyHeader(t *testing.T) {
	t.Run("full", func(t *testing.T) {
		g := NewWithT(t)

		a := http.Header{}
		a.Set("a", "b")
		a.Add("a", "c")
		a.Add("d", "e")
		b := http.Header{}
		c := a.Clone()
		copyHeader(a, b)

		g.Expect(a).ShouldNot(BeEmpty())
		g.Expect(b).Should(Equal(c))
		g.Expect(b).Should(HaveLen(2))
		g.Expect(b["A"]).Should(HaveLen(2))
		g.Expect(b["D"]).Should(HaveLen(1))
	})

	t.Run("empty", func(t *testing.T) {
		g := NewWithT(t)

		a := http.Header{}
		b := http.Header{}
		copyHeader(a, b)

		g.Expect(b).Should(BeEmpty())
	})

	t.Run("self", func(t *testing.T) {
		g := NewWithT(t)

		a := http.Header{}
		a.Add("a", "b")
		copyHeader(a, a)

		g.Expect(a).Should(HaveLen(1))
		g.Expect(a["A"]).Should(HaveLen(1), "keys were not duplicated onto themselves")
	})
}

func TestProxyUrl(t *testing.T) {
	t.Run("no mutation of base", func(t *testing.T) {
		g := NewWithT(t)
		base := v2.DefaultDnsdbServer.String()
		proxyURL(v2.DefaultDnsdbServer, mustParseURL("http://localhost/dnsdbfront/dnsdb/v2/ping"))
		g.Expect(v2.DefaultDnsdbServer.String()).Should(Equal(base))
	})

	t.Run("no params", func(t *testing.T) {
		g := NewWithT(t)

		a := proxyURL(v2.DefaultDnsdbServer, mustParseURL("http://localhost/dnsdbfront/dnsdb/v2/ping"))
		g.Expect(a.String()).Should(Equal("https://api.dnsdb.info/dnsdb/v2/ping"))
	})

	t.Run("with params", func(t *testing.T) {
		g := NewWithT(t)

		a := proxyURL(v2.DefaultDnsdbServer, mustParseURL("http://localhost/dnsdbfront/dnsdb/v2/ping?b=2&a=1"))
		g.Expect(a.String()).Should(Equal("https://api.dnsdb.info/dnsdb/v2/ping?a=1&b=2"), "params copied and normalized")
	})

	t.Run("out of path", func(t *testing.T) {
		g := NewWithT(t)

		a := proxyURL(v2.DefaultDnsdbServer, mustParseURL("http://localhost/foo/bar"))
		g.Expect(a.String()).Should(Equal("https://api.dnsdb.info/foo/bar"))
	})
}

func TestProxyRequest(t *testing.T) {
	g := NewWithT(t)

	remoteAddr := "1.2.3.4"
	apikey := "fooapikey"

	req, err := http.NewRequest(http.MethodGet, "http://localhost/dnsdbfront/dnsdb/v2/ping", nil)
	g.Expect(err).ShouldNot(HaveOccurred())
	ctx := context.TODO()
	req = req.WithContext(ctx)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", apikey)
	req.RemoteAddr = remoteAddr

	proxied, err := proxyRequest(req)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(proxied.URL.String()).Should(Equal("https://api.dnsdb.info/dnsdb/v2/ping"), "url is translated")
	g.Expect(proxied.Header.Get("Accept")).Should(Equal("application/x-ndjson"), "accept header is altered")
	g.Expect(proxied.Header.Get("X-API-Key")).Should(Equal(apikey), "apikey is copied over")
	g.Expect(proxied.Header.Get("X-Forwarded-For")).Should(Equal(remoteAddr), "x-forwarded-for is set")
	g.Expect(proxied.Context()).Should(Equal(ctx), "context is copied over")
}

func mustParseURL(s string) *url.URL {
	if u, err := url.Parse(s); err != nil {
		panic(err)
	} else {
		return u
	}
}
