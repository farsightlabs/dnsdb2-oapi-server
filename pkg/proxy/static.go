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
	_ "embed"
	"net/http"

	"github.com/elnormous/contenttype"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const (
	textHtml    = "text/html"
	title       = "OpenAPI-Compatible DNSDB Proxy"
	icon        = "https://www.farsightsecurity.com/assets/media/png/favicon-152.png"
	css         = ""
	renderFlags = html.CompletePage
)

var mediaTypes = []contenttype.MediaType{
	contenttype.NewMediaType(textHtml),
}

//go:embed api.yaml
var apiText []byte

//go:embed README.md
var readmeText []byte

// rendered html README
var readmeHtml []byte

func init() {
	readmeHtml = markdown.ToHTML(readmeText, parser.New(),
		html.NewRenderer(html.RendererOptions{
			Title: title,
			CSS:   css,
			Icon:  icon,
			Flags: renderFlags,
		}),
	)
}

func InfoHandler(w http.ResponseWriter, req *http.Request) {
	accepted, _, err := contenttype.GetAcceptableMediaType(req, mediaTypes)

	if err == nil {
		switch accepted.String() {
		case textHtml:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(readmeHtml)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(readmeText)
}

func ApiHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(apiText)
}
