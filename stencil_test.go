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
	//"io"
	//"io/ioutil"
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

	const (
		templateFile = "./testdata/template.html"
		siteRoot     = "./testdata"
	)

	input := `

			root ` + siteRoot + `
			stencil / {
				template ` + templateFile + `
			}
		`

	// fmt.Println(input)

	c := caddy.NewTestController("http", input)
	err := stencil.Setup(c)
	if err != nil {
		t.Fatalf("Something went wrong loading the controller: %v\n", err)
	}

	mids := httpserver.GetConfig(c).Middleware()
	handler := mids[0](httpserver.EmptyNext).(stencil.Stencil)

	// fmt.Printf("handler: %v\n", handler.Configs[0])
	// fmt.Printf("TemplateFiles: %v\n", handler.Configs[0].Template)
	// fmt.Printf("Next: %v\n", handler.Next)
	handler.Next = staticfiles.FileServer{Root: http.Dir(siteRoot)}

	req, err := http.NewRequest("GET", "/index.html", nil)
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}

	req = req.WithContext(context.WithValue(req.Context(), httpserver.OriginalURLCtxKey, *req.URL))

	// fmt.Println(req)

	rec := httptest.NewRecorder()
	_, err = handler.ServeHTTP(rec, req)
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println(rec.Code)
	// fmt.Printf("%v\n", rec.Header().Get("Content-Type"))
	respBody := rec.Body.String()
	fmt.Printf("%v\n", respBody)

}
