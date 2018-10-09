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
// Adapted from the Caddy markdown plugin by Light Code Labs, LLC.
// https://github.com/mholt/caddy/tree/master/caddyhttp/markdown
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
)

// NoneParser is the parser for plaintext with no metadata.
type NoneParser struct {
	metadata Metadata
	body     *bytes.Buffer
}

// Type returns the kind of parser this struct is.
func (n *NoneParser) Type() string {
	return "None"
}

// Parse prepases and parses the metadata and body
func (n *NoneParser) Parse(b []byte) bool {
	m := make(map[string]interface{})
	n.metadata = NewMetadata(m)
	n.body = bytes.NewBuffer(b)

	return true
}

// Metadata returns parsed metadata.  It should be called
// only after a call to Parse returns without error.
func (n *NoneParser) Metadata() Metadata {
	return n.metadata
}

// Body returns parsed body.  It should be called
// only after a call to Parse returns without error.
func (n *NoneParser) Body() []byte {
	return n.body.Bytes()
}
