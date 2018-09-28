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
	"testing"
	"text/template"

	"github.com/jimjimovich/caddy-stencil"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func TestSetup(t *testing.T) {
	var tests = []struct {
		inputCongig    string
		shouldErr      bool
		expectedConfig []stencil.Config
	}{
		// Config with only defaults
		{
			"stencil /",
			false,
			[]stencil.Config{{
				PathScope: "/",
				Extensions: map[string]struct{}{
					".html": {},
					".json": {},
				},
				Template:      stencil.GetDefaultTemplate(),
				TemplateFiles: make(map[string]*stencil.CachedFileInfo),
			}}},
		// Config with named template, ext configured, and multiple entries
		{
			`
				stencil /blog
				stensil /test {
					template test ./testdata/index.html
					ext .txt .html
				}
			`,
			false,
			[]stencil.Config{{
				PathScope: "/blog",
				Extensions: map[string]struct{}{
					".html": {},
					".json": {},
				},
				Template:      stencil.GetDefaultTemplate(),
				TemplateFiles: make(map[string]*stencil.CachedFileInfo),
			},
				{
					PathScope: "/test",
					Extensions: map[string]struct{}{
						".txt":  {},
						".html": {},
					},
					Template: buildTemplate(map[string]string{
						"test": "testdata/index.html",
					}),
					TemplateFiles: map[string]*stencil.CachedFileInfo{
						"test": &stencil.CachedFileInfo{"testdata/index.html", nil},
					},
				}}},
	}

	for i, test := range tests {
		c := caddy.NewTestController("http", test.inputCongig)
		err := stencil.Setup(c)
		if err == nil && test.shouldErr {
			t.Errorf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d errored, but it shouldn't have; got '%v'", i, err)
		}

		mids := httpserver.GetConfig(c).Middleware()
		if len(mids) == 0 {
			t.Fatal("Expected middleware, got 0 instead")
		}

		handler := mids[0](httpserver.EmptyNext)
		myHandler, ok := handler.(stencil.Stencil)

		if !ok {
			t.Fatalf("Expected handler to be type Stencil, got: %#v", handler)
		}

		for j, singleConfig := range myHandler.Configs {

			// // Test correct path is set
			if singleConfig.PathScope != test.expectedConfig[j].PathScope {
				t.Errorf("Expected %v PathScope, but got %v", test.expectedConfig[j].PathScope, singleConfig.PathScope)
			}

			// Test extensions
			if len(test.expectedConfig[j].Extensions) != len(singleConfig.Extensions) {
				t.Errorf("Expected %v extensions, got: %v", len(test.expectedConfig[j].Extensions), len(singleConfig.Extensions))
			}

			for v, _ := range test.expectedConfig[j].Extensions {
				if _, ok := singleConfig.Extensions[v]; !ok {
					t.Errorf("Expected extensions to contain %v", v)
				}
			}

			// // Test template files
			if len(singleConfig.TemplateFiles) != len(test.expectedConfig[j].TemplateFiles) {
				t.Errorf("Expected %v TemplateFiles, got: %v", len(test.expectedConfig[j].TemplateFiles), len(singleConfig.TemplateFiles))
			}

			// Test TemplateFile Paths
			for tfk, tf := range test.expectedConfig[j].TemplateFiles {
				if tf.Path != singleConfig.TemplateFiles[tfk].Path {
					t.Errorf("Expected TemplateFile Path of %v , got: %v", tf.Path, singleConfig.TemplateFiles[tfk].Path)
				}
			}

			// Test tempate (and/or default) was loaded
			realTemplateNames := collectTemplateNames(singleConfig.Template.Templates())
			expectedTemplateNames := collectTemplateNames(test.expectedConfig[j].Template.Templates())

			for name, _ := range expectedTemplateNames {
				if _, ok := realTemplateNames[name]; !ok {
					t.Errorf("Expected template names to contain %v, but got %v", name, realTemplateNames[name])
				}
			}

			// Test that we have the expected amount of templates
			if len(test.expectedConfig[j].Template.Templates()) != len(singleConfig.Template.Templates()) {
				t.Errorf("Expected %v default templates, got: %v", len(test.expectedConfig[j].Template.Templates()), len(singleConfig.Template.Templates()))
			}
		}
	}
}

func buildTemplate(templates map[string]string) *template.Template {
	t := stencil.GetDefaultTemplate()
	for k, v := range templates {
		stencil.SetTemplate(t, k, v)
	}
	return t
}

func collectTemplateNames(templates []*template.Template) map[string]interface{} {
	templateNames := make(map[string]interface{})
	var i interface{}
	for _, template := range templates {
		name := template.Name()
		templateNames[name] = i
	}
	return templateNames
}
