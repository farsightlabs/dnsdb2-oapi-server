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
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes(Prefix) {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func routes(prefix string) Routes {
	prefix = normalizePrefix(prefix)

	return Routes{
		Route{
			"Info",
			http.MethodGet,
			fmt.Sprintf("%s/", prefix),
			InfoHandler,
		},

		Route{
			"API",
			http.MethodGet,
			fmt.Sprintf("%s/api.yaml", prefix),
			ApiHandler,
		},

		Route{
			"Lookup-RData",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/lookup/rdata/{type}/{value}", prefix),
			ProxySafHandler,
		},

		Route{
			"Lookup-RData-RRtype",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/lookup/rdata/{type}/{value}/{rrtype}", prefix),
			ProxySafHandler,
		},

		Route{
			"Lookup-RRSet",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/lookup/rrset/{type}/{value}", prefix),
			ProxySafHandler,
		},

		Route{
			"Lookup-RRSet-Bailiwick",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/lookup/rrset/{type}/{value}/{rrtype}/{bailiwick}", prefix),
			ProxySafHandler,
		},

		Route{
			"Lookup-RRSet-RRType",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/lookup/rrset/{type}/{value}/{rrtype}", prefix),
			ProxySafHandler,
		},

		Route{
			"Flex",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/{method}/{key}/{value}", prefix),
			ProxySafHandler,
		},

		Route{
			"Ping",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/ping", prefix),
			ProxyPassHandler,
		},

		Route{
			"RateLimit",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/rate_limit", prefix),
			ProxyPassHandler,
		},

		Route{
			"Summarize-RData",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/summarize/rdata/{type}/{value}", prefix),
			ProxySafHandler,
		},

		Route{
			"Summarize-RData-RRType",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/summarize/rdata/{type}/{value}/{rrtype}", prefix),
			ProxySafHandler,
		},

		Route{
			"Summarize-RRSet",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/summarize/rrset/{type}/{value}", prefix),
			ProxySafHandler,
		},

		Route{
			"Summarize-RRSet-Bailiwick",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/summarize/rrset/{type}/{value}/{rrtype}/{bailiwick}", prefix),
			ProxySafHandler,
		},

		Route{
			"Summarize-RRSet-RRType",
			http.MethodGet,
			fmt.Sprintf("%s/dnsdb/v2/summarize/rrset/{type}/{value}/{rrtype}", prefix),
			ProxySafHandler,
		},
	}
}
