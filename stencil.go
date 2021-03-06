// Copyright 2018 Jim Mendenhall
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Adapted from various Caddy plugins by Light Code Labs, LLC.
// https://github.com/mholt/caddy
// Significant modifications have been made.
//
// Original License
//
// Copyright 2015 Light Code Labs, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stencil

import (
	"bytes"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type Stencil struct {
	// Server root
	Root string

	// Jail the requests to site root with a mock file system
	FileSys http.FileSystem

	// Next HTTP handler in the chain
	Next httpserver.Handler

	// The list of stencil configurations
	Configs []*Config

	BufPool *sync.Pool
}

// Config stores stencil middleware configurations.
type Config struct {
	// Base path to match
	PathScope string

	// List of extensions to consider as stencil files
	Extensions map[string]struct{}

	// Template(s) to render with
	Template *template.Template

	// a pair of template's name and its underlying file information
	TemplateFiles map[string]*CachedFileInfo
}

type CachedFileInfo struct {
	Path string
	Fi   os.FileInfo
}

// ServeHTTP implements the http.Handler interface.
func (st Stencil) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	var cfg *Config
	for _, c := range st.Configs {
		if httpserver.Path(r.URL.Path).Matches(c.PathScope) {
			cfg = c
			break
		}
	}
	if cfg == nil {
		return st.Next.ServeHTTP(w, r)
	}

	originalMethod := r.Method
	// If HEAD request, temporarily set to GET so that staticfiles or proxy
	// will send content and we can calculate content-length correctly for HEAD requests
	if r.Method == http.MethodHead {
		r.Method = http.MethodGet
	}

	fpath := r.URL.Path

	// get a buffer from the pool and make a response recorder
	buf := st.BufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer st.BufPool.Put(buf)

	// only buffer the response when we want to execute a stencil
	shouldBuf := func(status int, header http.Header) bool {
		// see if this request matches a stencil extension
		reqExt := path.Ext(fpath)
		for ext := range cfg.Extensions {
			// do not buffer if redirect or error
			if status >= 300 {
				return false
			}
			if reqExt == "" {
				// request has no extension, so check response Content-Type
				ct := mime.TypeByExtension(ext)
				if ct != "" && strings.Contains(ct, header.Get("Content-Type")) {
					return true
				}
			} else if reqExt == ext {
				return true
			}
		}
		return false
	}

	// prepare a buffer to hold the response, if applicable
	rb := httpserver.NewResponseBuffer(buf, w, shouldBuf)

	// pass request up the chain to let another middleware provide us content
	// this will most likely come from staticfiles or proxy
	code, err := st.Next.ServeHTTP(rb, r)
	if !rb.Buffered() || code >= 300 || err != nil {
		return code, err
	}

	// create an execution context
	ctx := httpserver.NewContextWithHeader(w.Header())
	ctx.Root = st.FileSys
	ctx.Req = r
	ctx.URL = r.URL

	html, err := cfg.Stencil(title(fpath), rb.Buffer, ctx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// reset to original HTTP method if we changed it
	if r.Method != originalMethod {
		r.Method = originalMethod
	}

	// copy the buffered header into the real ResponseWriter
	rb.CopyHeader()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	lastModTime, _ := time.Parse(http.TimeFormat, w.Header().Get("Last-Modified"))
	http.ServeContent(rb.StatusCodeWriter(w), r, fpath, lastModTime, bytes.NewReader(html))

	return 0, nil
}

// title gives a backup generated title for a page
func title(p string) string {
	return strings.TrimSuffix(path.Base(p), path.Ext(p))
}
