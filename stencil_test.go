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
// Based on tests from Caddy Server that are Copyright 2015 Light Code Labs, LLC
// https://github.com/mholt/caddy

package stencil_test

import (
	"context"
	"fmt"
	"strings"
	//"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	//"text/template"
	"github.com/jimjimovich/caddy-stencil"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/mholt/caddy/caddyhttp/staticfiles"
)

func TestStencil(t *testing.T) {

	tests := []struct {
		siteRoot     string
		inputConfig  string
		getPath      string
		expectedFile string
	}{
		{
			"./testdata/json",
			`stencil / {
				ext .json
				template ./testdata/json/template.html
			}
			`,
			"/44418.json",
			"/44418_expected.html",
		},
		{
			"./testdata/json",
			`stencil / {
				ext .json
				template ./testdata/json/template.html
			}
			`,
			"/search.json",
			"/search_expected.html",
		},
		{
			"./testdata/json_frontmatter",
			`stencil / {
				template ./testdata/json_frontmatter/template.html
			}
			`,
			"/index.html",
			"/index_expected.html",
		},
		{
			"./testdata/yaml_frontmatter",
			`stencil / {
				template ./testdata/yaml_frontmatter/template.html
			}
			`,
			"/index.html",
			"/index_expected.html",
		},
		{
			"./testdata/toml_frontmatter",
			`stencil / {
				template ./testdata/toml_frontmatter/template.html
			}
			`,
			"/index.html",
			"/index_expected.html",
		},
		{
			"./testdata/problem_files",
			`stencil / {
				template ./testdata/problem_files/template.html
			}
			`,
			"/badjson.html",
			"/badjson_expected.html",
		},
		{
			"./testdata/problem_files",
			`stencil / {
				template ./testdata/problem_files/template.html
			}
			`,
			"/badyaml.html",
			"/badyaml_expected.html",
		},
		{
			"./testdata/json_frontmatter",
			`stencil / {
				template ./testdata/json_frontmatter/template.html
				template alternate ./testdata/json_frontmatter/alternate.html
			}
			`,
			"/alt.html",
			"/alt_expected.html",
		},
	}

	for _, test := range tests {
		c := caddy.NewTestController("http", test.inputConfig)
		err := stencil.Setup(c)
		if err != nil {
			t.Fatalf("Something went wrong loading the controller: %v\n", err)
		}

		mids := httpserver.GetConfig(c).Middleware()
		handler := mids[0](httpserver.EmptyNext).(stencil.Stencil)
		handler.Next = staticfiles.FileServer{Root: http.Dir(test.siteRoot)}

		req, err := http.NewRequest("GET", test.getPath, nil)
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}

		req = req.WithContext(context.WithValue(req.Context(), httpserver.OriginalURLCtxKey, *req.URL))

		rec := httptest.NewRecorder()
		_, err = handler.ServeHTTP(rec, req)
		if err != nil {
			t.Fatal(err)
		}

		respBody := rec.Body
		expectedBody := expected(test.siteRoot + test.expectedFile)

		got := strings.TrimSpace(respBody.String())
		expected := strings.TrimSpace(string(expectedBody))

		if got != expected {
			t.Fatalf("Expected:\n %v\n Got:\n %v\n", expected, got)
		}
	}

}

func expected(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}

	return b
}
