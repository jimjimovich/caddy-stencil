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
	"bytes"
	"encoding/json"
)

// ValidJSONParser is the Parser for already validated JSON
type ValidJSONParser struct {
	metadata Metadata
	body     *bytes.Buffer
}

// Type returns the kind of metadata parser implemented by this struct.
func (j *ValidJSONParser) Type() string {
	return "ValidJSON"
}

// Parse processes JSON and adds it to j.metadata. j.body is set to nil.
func (j *ValidJSONParser) Parse(by []byte) bool {
	// Figure out if this starts with an [ or { to see if an array
	var isArray = false
	
	if by[0] == bytes.TrimSpace([]byte("["))[0] {
		isArray = true
	}
	
	r := bytes.NewBuffer(by)
	
	var arrayData interface{}
	var data map[string]interface{}

	d := json.NewDecoder(r)
	d.UseNumber()

	if isArray {
		if err := d.Decode(&arrayData); err != nil {
			// there really shouldn't be an error here
		}
	} else {
		if err := d.Decode(&data); err != nil {
			// there really shouldn't be an error here
		}
	}

	if isArray {
		metaMap := make(map[string]interface{})
		metaMap["data"] = arrayData
		mdata := NewMetadata(metaMap)
		j.metadata = mdata
		j.body = bytes.NewBuffer(nil)
		return true
	} else {
		metaMap := make(map[string]interface{})
		metaMap["data"] = data
		mdata := NewMetadata(metaMap)
		j.metadata = mdata
		j.body = bytes.NewBuffer(nil)
		return true
	}
	
	return false
}

// Metadata returns parsed metadata.  It should be called
// only after a call to Parse returns without error.
func (j *ValidJSONParser) Metadata() Metadata {
	return j.metadata
}

// Body returns the body text.  It should be called only after a call to Parse returns without error.
func (j *ValidJSONParser) Body() []byte {
	return j.body.Bytes()
}
