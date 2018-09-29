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
// Adapted from the Caddy Markdown plugin by Light Code Labs, LLC.
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

package metadata

import (
	"bytes"

	"github.com/naoina/toml"
)

// TOMLParser is the Parser for TOML
type TOMLParser struct {
	metadata Metadata
	body     *bytes.Buffer
}

// Type returns the kind of parser this struct is.
func (t *TOMLParser) Type() string {
	return "TOML"
}

// Parse prepares and parses the metadata and body
func (t *TOMLParser) Parse(by []byte) bool {
	b := bytes.NewBuffer(by)
	meta, data := splitBuffer(b, "+++")
	if meta == nil || data == nil {
		return false
	}
	t.body = data

	m := make(map[string]interface{})
	if err := toml.Unmarshal(meta.Bytes(), &m); err != nil {
		return false
	}
	t.metadata = NewMetadata(m)

	return true
}

// Metadata returns parsed metadata.  It should be called
// only after a call to Parse returns without error.
func (t *TOMLParser) Metadata() Metadata {
	return t.metadata
}

// Body returns parser the body.  It should be called only after a call to Parse returns without error.
func (t *TOMLParser) Body() []byte {
	return t.body.Bytes()
}
