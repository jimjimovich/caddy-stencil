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
	"bytes"
	"encoding/json"
)

// JSONParser is the MetadataParser for JSON
type JSONParser struct {
	metadata Metadata
	body     *bytes.Buffer
}

// Type returns the kind of metadata parser implemented by this struct.
func (j *JSONParser) Type() string {
	return "JSON"
}

// Parse processes the document and prepares the metadata and body
func (j *JSONParser) Parse(by []byte) bool {
	// Valid JSON arrays may appear in JSON APIs, so we need to deal with them.
	// JSON arrays should not be used as front matter wihtout being wrapped
	// in a JSON object { }.  Any non-valid JSON arrays will not be processed
	// as JSON. The body object for valid JSON arrays is always returned as nil.
	//
	// Figure out if this starts with an [ or { to see if an array.
	var isArray = false
	if bytes.TrimSpace(by)[0] == []byte("[")[0] {
		isArray = true
	}

	var arrayData interface{}
	var data map[string]interface{}

	buf := bytes.NewBuffer(by)

	// If we have a JSON array, we should have valid JSON from an API with no body
	if isArray {
		err := json.Unmarshal(buf.Bytes(), &arrayData)
		if err != nil {
			return false
		}
		metaMap := make(map[string]interface{})
		metaMap["data"] = arrayData
		mdata := NewMetadata(metaMap)
		j.metadata = mdata
		j.body = bytes.NewBuffer(nil)
		return true
	} else {
		// Starts with "{", may be JSON document or another document with JSON
		// front matter. If valid JSON with no body, body is returned as nil.
		err := json.Unmarshal(buf.Bytes(), &data)
		if err != nil {
			var offset int

			jerr, ok := err.(*json.SyntaxError)
			if !ok {
				return false
			}

			offset = int(jerr.Offset)

			err = json.Unmarshal(buf.Next(offset-1), &data)
			if err != nil {
				return false
			}

			j.body = bytes.NewBuffer(buf.Bytes())

		} else {
			// There was no error processing the entire document, so we have
			// valid JSON. We set the body to nil.
			j.body = bytes.NewBuffer(nil)
		}

		metaMap := make(map[string]interface{})
		metaMap["data"] = data
		mdata := NewMetadata(metaMap)
		j.metadata = mdata

		return true
	}
}

// Metadata returns parsed metadata.  It should be called
// only after a call to Parse returns without error.
func (j *JSONParser) Metadata() Metadata {
	return j.metadata
}

// Body returns the body text.  It should be called only after a call to Parse returns without error.
func (j *JSONParser) Body() []byte {
	return j.body.Bytes()
}
