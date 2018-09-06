package templates

import (
	"log"
	"testing"

	"gopkg.in/yaml.v2"
)

var input = make(map[string]interface{})
var data = `
a: Easy!
b: 
  c: '2'
  h: 'ewff'
`

func init() {
	err := yaml.Unmarshal([]byte(data), input)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func TestTmplStrConst(t *testing.T) {
	input := "log.console.as.json"
	output := tmplStrConst(input)
	expected := "LOG_CONSOLE_AS_JSON"

	if output != expected {
		t.Errorf("strCon conversion was incorrect, got: %s, want: %s", output, expected)
	}
}

func TestTmplStr(t *testing.T) {
	output := tmplString("b.c", input)

	if output != "2" {
		t.Errorf("string conversion was incorrect, got: %s, want: %v", output, "2")
	}
}

func TestTmpObjArray(t *testing.T) {
	output := tmplObjectArray("b", input)

	if len(output) == 2 {
		if !(output[0].Key == "c" && output[0].Value == "2" && output[1].Key == "h" && output[1].Value == "ewff") {
			t.Errorf("ObjArray didn't match expected values, got: [{%s, %s}, {%s, %s}], want: [{c, 2}, {h, ewff}]", output[0].Key, output[0].Value, output[1].Key, output[1].Value)
		}
	} else {
		t.Errorf("ObjArray didn't match length, got: %s, want: %v", output, 2)
	}
}
