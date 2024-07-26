Go HTML Template Autoescaping
July 26, 2024

Go's [`html/template`](https://pkg.go.dev/html/template) library does contextual autoescaping. The package tries to ensure any input data is safe to insert into the HTML. It transforms input data differently depending on whether you're inserting it as straight up text in the HTML, or JavaScript code inside a `<script>` tag.

It turns out you can see what autoescaping `html/template` is doing under the hood. See the code below for a few examples.

One thing I thought was interesting is that Go tries to detect whether to JavaScript escape the input [based on the mime type of the the `<script>` tag](https://github.com/golang/go/blob/d8c7230c97ca5639389917cc235175bfe2dc50ab/src/html/template/js.go#L450-L485).

```go
package main

import (
	"bytes"
	"html/template"
	"testing"
)

func TestRunTemplate(t *testing.T) {
	type testData struct {
		inputTemplate          string
		expectedOutputTemplate string
		input                  any
		expectedOutput         any
	}

	runTest := func(test testData) {
		t.Helper()
		outputTemplate, output := runTemplate(t, test.inputTemplate, test.input)
		if outputTemplate != test.expectedOutputTemplate {
			t.Errorf("outputTemplate: got %s, want %s\n", outputTemplate, test.expectedOutputTemplate)
		}
		if output != test.expectedOutput {
			t.Errorf("output: got %s, want %s\n", output, test.expectedOutput)
		}
	}

	runTest(testData{
		inputTemplate:          `<script>{{.}}</script>`,
		expectedOutputTemplate: `<script>{{. | _html_template_jsvalescaper}}</script>`,
		input:                  `<body>`,
		expectedOutput:         `<script>"\u003cbody\u003e"</script>`,
	})
	runTest(testData{
		inputTemplate:          `<script type="foo">{{.}}</script>`,
		expectedOutputTemplate: `<script type="foo">{{. | _html_template_htmlescaper}}</script>`,
		input:                  `<body>`,
		expectedOutput:         `<script type="foo">&lt;body&gt;</script>`,
	})
	runTest(testData{
		inputTemplate:          `<script type="bar">{{.}}</script>`,
		expectedOutputTemplate: `<script type="bar">{{. | _html_template_htmlescaper}}</script>`,
		input:                  template.HTML(`<body>`),
		expectedOutput:         `<script type="bar"><body></script>`,
	})
}

func runTemplate(t *testing.T, templateText string, input any) (outputTemplate, output string) {
	tmpl := template.Must(template.New("").Parse(templateText))

	var out bytes.Buffer
	if err := tmpl.Execute(&out, input); err != nil {
		t.Fatalf("template.Execute: %v", err)
	}

	return tmpl.Tree.Root.String(), out.String()
}
```
