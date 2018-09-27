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

package stencil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// FileInfo represents a file in a particular server context. It wraps the os.FileInfo struct.
type FileInfo struct {
	os.FileInfo
	ctx httpserver.Context
}

// Stencil processes the contents of a page in b. It parses the metadata
// (if any) and uses the template (if found).
func (c *Config) Stencil(title string, body io.Reader, ctx httpserver.Context) ([]byte, error) {
	mdata, err := parseBody(body, false)
	if err != nil {
		// If the error is because of a JSON array, retry to process as array
		if strings.Contains(err.Error(), "cannot unmarshal array") {
			mdata, _ := parseBody(body, true)
			return execTemplate(c, mdata, ctx)
		}
	}

	return execTemplate(c, mdata, ctx)
}

func parseBody(body io.Reader, array bool) (Metadata, error) {
	d := json.NewDecoder(body)
	d.UseNumber()

	var b map[string]interface{}
	var a interface{}

	if array {
		if err := d.Decode(&a); err != nil {
			return Metadata{}, err
		}
	} else {
		if err := d.Decode(&b); err != nil {
			// If invalid character, put whole body into the body variable
			if strings.Contains(err.Error(), "invalid character") {
				mdata := Metadata{
					Variables: make(map[string]interface{}),
				}
				buf := new(bytes.Buffer)
				buf.ReadFrom(d.Buffered())
				mdataBody := buf.String()
				mdata.Variables["body"] = mdataBody
				return mdata, nil
			}
			return Metadata{}, err
		}
	}

	// No error decoding, so we have JSON or JSON + body
	if array {
		metaMap := make(map[string]interface{})
		metaMap["data"] = a
		mdata := NewMetadata(metaMap)
		if d.More() {
			buf := new(bytes.Buffer)
			buf.ReadFrom(d.Buffered())
			mdataBody := buf.String()
			mdata.Variables["body"] = mdataBody
		}
		return mdata, nil
	} else {
		mdata := NewMetadata(b)
		if d.More() {
			// var buf bytes.Buffer

			buf := new(bytes.Buffer)
			buf.ReadFrom(d.Buffered())
			mdataBody := buf.String()
			fmt.Println(mdataBody)

			mdata.Variables["body"] = mdataBody
		}
		return mdata, nil
	}
}
