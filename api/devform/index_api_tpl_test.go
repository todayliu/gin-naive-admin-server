package devform

import (
	"strings"
	"testing"
)

func TestIndexAPITsTemplate_QueryReturnTypeHasNoTripleClose(t *testing.T) {
	data := &genTemplateData{
		EntityName:  "DemoDict",
		RouteGroup:  "demo_dict",
		Description: "",
		Fields: []GenField{
			{JSONName: "demoDictCode", TSType: "string"},
		},
		QueryFieldsTS: []tsFieldVue{
			{JSONName: "demoDictCode", TSType: "string", Comment: "c"},
		},
	}
	b, err := execTpl("index.api.ts.tpl", data)
	if err != nil {
		t.Fatal(err)
	}
	out := string(b)
	if !strings.Contains(out, "static query") {
		t.Fatalf("missing query method:\n%s", out)
	}
	// list() 合法地会有 DemoDictRow>>>（三层泛型）；query() 不应在同一行写成 BaseResponse<DemoRow>>>
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "static query") && strings.Contains(line, ">>>") {
			t.Fatalf("query() return type should not use single-line triple >:\n%s", line)
		}
	}
}
