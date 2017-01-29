// Copyright 2012 Arne Roomann-Kurrik
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

package json

import (
	"fmt"
	"reflect"
	"testing"
)

type TestCase struct {
	Raw    string
	Result interface{}
}

var errors = map[string]TestCase{
	"HTML": {
		Raw:    "<!DOCTYPE html><html><body>Foo</body></html>",
		Result: "Unrecognized type in ' --><<-- !DOCTYPE '",
	},
	"Blank": {
		Raw:    "",
		Result: "Unrecognized type in ' --><-- '",
	},
	"Empty": {
		Raw:    "    ",
		Result: "Unrecognized type in '    --><-- '",
	},
}

var cases = map[string]TestCase{
	"Number": {
		Raw:    "1234",
		Result: int64(1234),
	},
	"Number - negative": {
		Raw:    "-1234",
		Result: int64(-1234),
	},
	"Number - float": {
		Raw:    "1234.5678",
		Result: float64(1234.5678),
	},
	"Number - negative float": {
		Raw:    "-1234.5678",
		Result: float64(-1234.5678),
	},
	"String": {
		Raw:    "\"foobar\"",
		Result: "foobar",
	},
	"String with encoded UTF-8": {
		Raw:    "\"\\u6211\\u7231\\u4f60\"",
		Result: "我爱你",
	},
	"String with unencoded UTF-8": {
		Raw:    "\"我爱你\"",
		Result: "我爱你",
	},
	"String with big-U encoded multibyte UTF-8": {
		Raw:    "\"\\U0001D11E\"",
		Result: "𝄞",
	},
	"String with octal encoded multibyte UTF-8": {
		Raw:    "\"\\360\\235\\204\\236\"",
		Result: "𝄞",
	},
	"String with hex encoded multibyte UTF-8": {
		Raw:    "\"\\xF0\\x9D\\x84\\x9E\"",
		Result: "𝄞",
	},
	"String with hex encoded single byte UTF-8": {
		Raw:    "\"\\xE2\\x9D\\xA4\"",
		Result: "❤",
	},
	"String with encoded UTF-8 and backslash": {
		Raw:    "\"10\\\\10 ~ \\u2764\"",
		Result: "10\\10 ~ ❤",
	},
	"Invalid string with small-U encoded multibyte UTF-8": {
		Raw:    "\"\\uD834\\uDD1E\"",
		Result: "𝄞",
	},
	"String with backslash": {
		Raw:    "\"10\\\\10\"",
		Result: "10\\10",
	},
	"String with backslash and tab": {
		Raw: "\"10\\\\\t10\"",
		Result: "10\\	10",
	},
	"String with backslash and backspace": {
		Raw:    "\"10\\\\\b10\"",
		Result: "10\\\b10",
	},
	"String with escaped forward slash": {
		Raw:    "\"\\\\\\/\"",
		Result: "\\/",
	},
	"String with just backslash": {
		Raw:    "\"\\\\\"",
		Result: "\\",
	},
	"String with encoded emoji": {
		Raw:    "\"EMOJI \\ud83d\\ude04 \\ud83d\\ude03 \\ud83d\\ude00 \\ud83d\\ude0a\"",
		Result: "EMOJI 😄 😃 😀 😊",
	},
	"Object": {
		Raw: "{\"foo\":\"bar\"}",
		Result: map[string]interface{}{
			"foo": "bar",
		},
	},
	"Object with spaces": {
		Raw: "{ \"foo\" : \"bar\" }",
		Result: map[string]interface{}{
			"foo": "bar",
		},
	},
	"Object with UTF-8 value": {
		Raw: "{ \"foo\" : \"\\u6211\" }",
		Result: map[string]interface{}{
			"foo": "我",
		},
	},
	"Object with tabs": {
		Raw: "{	\"foo\"	:	\"bar\"	}",
		Result: map[string]interface{}{
			"foo": "bar",
		},
	},
	"Object with empty nested object": {
		Raw: "{ \"foo\": {}}",
		Result: map[string]interface{}{
			"foo": map[string]interface{}{},
		},
	},
	"Object with empty nested array": {
		Raw: "{\"foo\": []}",
		Result: map[string]interface{}{
			"foo": []interface{}{},
		},
	},
	"Array": {
		Raw: "[1234,\"foobar\"]",
		Result: []interface{}{
			int64(1234),
			"foobar",
		},
	},
	"Array with spaces": {
		Raw: "[ 1234 , \"foobar\" ]",
		Result: []interface{}{
			int64(1234),
			"foobar",
		},
	},
	"Array with tabs": {
		Raw: "[	1234	,	\"foobar\"	]",
		Result: []interface{}{
			int64(1234),
			"foobar",
		},
	},
	"Array with multiple tabs": {
		Raw: "[				1234,\"foobar\"]",
		Result: []interface{}{
			int64(1234),
			"foobar",
		},
	},
	"Array with no contents": {
		Raw:    "[]",
		Result: []interface{}{},
	},
	"Array with empty object": {
		Raw: "[{}]",
		Result: []interface{}{
			map[string]interface{}{},
		},
	},
}

func TestCases(t *testing.T) {
	var (
		err    error
		decode interface{}
	)
	for desc, testcase := range cases {
		if err = Unmarshal([]byte(testcase.Raw), &decode); err != nil {
			t.Fatalf("Error decoding '%v': %v", desc, err)
		}
		if !reflect.DeepEqual(decode, testcase.Result) {
			t.Logf("%v\n", reflect.TypeOf(decode))
			t.Logf("%v\n", reflect.TypeOf(testcase.Result))
			if reflect.TypeOf(decode) == reflect.TypeOf("") {
				t.Logf("Decode: %v\n", []byte(decode.(string)))
				t.Logf("Expected: %v\n",
					[]byte(testcase.Result.(string)))
			}
			t.Fatalf("Problem decoding '%v' Expected: %v, Got %v",
				desc, testcase.Result, decode)
		}
	}
}

func TestErrors(t *testing.T) {
	var (
		err    error
		str    string
		res    string
		decode interface{}
	)
	for desc, tcase := range errors {
		if err = Unmarshal([]byte(tcase.Raw), &decode); err == nil {
			t.Fatalf("Expected error for '%v': %v", desc, tcase.Raw)
		}
		str = fmt.Sprintf("%v", err)
		res = tcase.Result.(string)
		if str != res {
			t.Fatalf("Invalid error '%v' expected '%v'", str, res)
		}
	}
}
