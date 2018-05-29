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
// Adapted from the Caddy Stencil plugin by Light Code Labs, LLC.
// Significant modifications have been made.
//
// Original License
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
	html     *bytes.Buffer
}

// Type returns the kind of metadata parser implemented by this struct.
func (j *JSONParser) Type() string {
	return "JSON"
}

// Init prepares the metadata metadata/html file and parses it
func (j *JSONParser) Init(b *bytes.Buffer) bool {
	m := make(map[string]interface{})

	err := json.Unmarshal(b.Bytes(), &m)
	if err != nil {
		var offset int

		jerr, ok := err.(*json.SyntaxError)
		if jerr == nil && !ok {
			// We have an error that is not a syntax error
			// Try to wrapp the json (probably an array) in { data: }
			sb := b.Bytes()
			sb = append([]byte(`{ "data": `), sb...)
			sb = append(sb, '}')
			err = json.Unmarshal(sb, &m)
			if err != nil {
				return false
			}
			j.metadata = NewMetadata(m)
			j.html = bytes.NewBuffer([]byte{})
			return true
		}

		if jerr != nil && !ok {
			// Seems like we don't have json
			return false
		}

		// If we got this far, We have a syntax error, which probably means
		// that we have data after our json, which should be added to .Doc.body
		offset = int(jerr.Offset)

		m = make(map[string]interface{})
		err = json.Unmarshal(b.Next(offset-1), &m)
		if err != nil {
			return false
		}
	}

	j.metadata = NewMetadata(m)
	j.html = bytes.NewBuffer(b.Bytes())

	return true
}

// Metadata returns parsed metadata.  It should be called
// only after a call to Parse returns without error.
func (j *JSONParser) Metadata() Metadata {
	return j.metadata
}

// Content returns the html body.  It should be called only after a call to Parse returns without error.
func (j *JSONParser) Content() []byte {
	return j.html.Bytes()
}
