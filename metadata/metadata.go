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

package metadata

import (
	"bufio"
	"bytes"
)

// Metadata stores a page's metadata
type Metadata struct {
	// Page title
	Title string

	// Page template name
	Template string

	// Variables to be used with Template
	Variables map[string]interface{}
}

// NewMetadata returns a new Metadata struct, loaded with the given map
func NewMetadata(parsedMap map[string]interface{}) Metadata {
	md := Metadata{
		Variables: make(map[string]interface{}),
	}
	md.load(parsedMap)

	return md
}

// load loads parsed values in parsedMap into Metadata
func (m *Metadata) load(parsedMap map[string]interface{}) {

	// Pull top level things out of data
	if data, ok := parsedMap["data"].(map[string]interface{}); ok {
		if title, ok := data["title"]; ok {
			m.Title, _ = title.(string)
		}

		// TODO: make template variable customizable in config
		if template, ok := data["template"]; ok {
			m.Template, _ = template.(string)
		}
	}

	m.Variables = parsedMap
}

// Parser is a an interface that must be satisfied by each parser
type Parser interface {
	// Initialize a parser
	Parse(b []byte) bool

	// Type of metadata
	Type() string

	// Parsed metadata.
	Metadata() Metadata

	// Raw body
	Body() []byte
}

// GetParser returns a parser for the given data
func GetParser(by []byte) Parser {
	// Get first characters of document to give us a clue about the type
	firstByte := bytes.TrimSpace(by)[0]

	switch firstByte {
	// First byte looks like JSON
	case []byte("{")[0], []byte("[")[0]:
		p := &JSONParser{}
		if p.Parse(by) {
			return p
		} else {
			n := &NoneParser{}
			if n.Parse(by) {
				return n
			}
		}

	// First byte looks like YAML
	case []byte("-")[0]:
		p := &YAMLParser{}
		if p.Parse(by) {
			return p
		} else {
			n := &NoneParser{}
			if n.Parse(by) {
				return n
			}
		}
	// First byte looks like TOML
	case []byte("+")[0]:
		p := &TOMLParser{}
		if p.Parse(by) {
			return p
		} else {
			n := &NoneParser{}
			if n.Parse(by) {
				return n
			}
		}

	// First byte looks like regular document
	default:
		n := &NoneParser{}
		if n.Parse(by) {
			return n
		}

		return nil
	}

	return nil
}

// parsers returns all available parsers
func parsers() []Parser {

	return []Parser{
		&TOMLParser{},
		&YAMLParser{},
		&JSONParser{},

		// This one must be last
		&NoneParser{},
	}
}

// Split out prefixed/suffixed metadata with given delimiter
func splitBuffer(b *bytes.Buffer, delim string) (*bytes.Buffer, *bytes.Buffer) {
	scanner := bufio.NewScanner(b)

	// Read and check first line
	if !scanner.Scan() {
		return nil, nil
	}
	if string(bytes.TrimSpace(scanner.Bytes())) != delim {
		return nil, nil
	}

	// Accumulate metadata, until delimiter
	meta := bytes.NewBuffer(nil)
	for scanner.Scan() {
		if string(bytes.TrimSpace(scanner.Bytes())) == delim {
			break
		}
		if _, err := meta.Write(scanner.Bytes()); err != nil {
			return nil, nil
		}
		if _, err := meta.WriteRune('\n'); err != nil {
			return nil, nil
		}
	}
	// Make sure we saw closing delimiter
	if string(bytes.TrimSpace(scanner.Bytes())) != delim {
		return nil, nil
	}

	// The rest is body
	body := new(bytes.Buffer)
	for scanner.Scan() {
		if _, err := body.Write(scanner.Bytes()); err != nil {
			return nil, nil
		}
		if _, err := body.WriteRune('\n'); err != nil {
			return nil, nil
		}
	}

	return meta, body
}
