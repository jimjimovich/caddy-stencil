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
	"io"
	"io/ioutil"
	"os"

	"github.com/jimjimovich/caddy-stencil/metadata"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// FileInfo represents a file in a particular server context. It wraps the os.FileInfo struct.
type FileInfo struct {
	os.FileInfo
	ctx httpserver.Context
}

// Stencil processes the contents of a page in r. It parses the metadata
// (if any) and uses the template (if found).
func (c *Config) Stencil(title string, r io.Reader, ctx httpserver.Context) ([]byte, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	parser := metadata.GetParser(contents)
	body := parser.Body()
	mdata := parser.Metadata()

	// set it as body for template
	mdata.Variables["body"] = string(body)

	// fixup title
	mdata.Variables["title"] = mdata.Title
	if mdata.Variables["title"] == "" {
		mdata.Variables["title"] = title
	}

	return execTemplate(c, mdata, ctx)
}
