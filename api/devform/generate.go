package devform

import (
	"archive/zip"
	"bytes"
	"embed"
	"fmt"
	"text/template"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed templates/*.tpl
var tmplFS embed.FS

// GenField 代码生成用字段描述
type GenField struct {
	ColumnName string
	JSONName   string
	GoName     string
	GoType     string
	StructTags string
	TSType     string
	ListShow   bool
	FormShow   bool
	IsQuery    bool
	QueryType  string
	DbType     string
	Comment    string
}

// genTemplateData 模板数据
type genTemplateData struct {
	PackageName   string
	PackageImport string
	TableName     string
	EntityName    string
	RouteGroup    string
	Description   string
	ServiceVar    string
	HasTimeType   bool
	Fields        []GenField
	QueryFields   []GenField
	QueryBlocks   []string
	QueryFieldsTS []tsFieldVue
	ListColumns   []listColVue
	FormFieldsVue []formFieldVue
	ViewPath      string
	ListTitle     string // 列表标题（说明或实体名，用于 pro-data-table）
	ShowSearch    bool   // 是否有查询条件（控制是否渲染 pro-search-form）
}

type tsFieldVue struct {
	JSONName string
	TSType   string
	Comment  string
}

type listColVue struct {
	Title    string
	JSONName string
	Ellipsis bool
}

type formFieldVue struct {
	Label      string
	ControlTpl string
}

func snakeToCamelJSON(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		low := strings.ToLower(parts[i])
		if i == 0 {
			parts[i] = low
		} else {
			parts[i] = strings.ToUpper(low[:1]) + low[1:]
		}
	}
	return strings.Join(parts, "")
}

func snakeToGoExported(s string) string {
	cc := snakeToCamelJSON(s)
	if cc == "" {
		return ""
	}
	return strings.ToUpper(cc[:1]) + cc[1:]
}

func serviceVar(entity string) string {
	if entity == "" {
		return "entity"
	}
	return strings.ToLower(entity[:1]) + entity[1:]
}

func buildGenFields(form *SysOnlineForm, fields []SysOnlineFormField) ([]GenField, error) {
	sort.Slice(fields, func(i, j int) bool {
		if fields[i].Sort != fields[j].Sort {
			return fields[i].Sort < fields[j].Sort
		}
		return fields[i].ID < fields[j].ID
	})
	out := make([]GenField, 0, len(fields))
	for _, f := range fields {
		if _, skip := reservedCol[strings.ToLower(f.ColumnName)]; skip {
			continue
		}
		if !validateIdent(f.ColumnName) {
			return nil, fmt.Errorf("非法列名: %s", f.ColumnName)
		}
		if !validateDBType(f.DbType) {
			return nil, fmt.Errorf("非法 db_type: %s", f.DbType)
		}
		gj := snakeToCamelJSON(f.ColumnName)
		gn := snakeToGoExported(f.ColumnName)
		dt := strings.ToLower(f.DbType)
		var goType, tsType, gormType string
		switch dt {
		case "varchar":
			L := f.Length
			if L <= 0 || L > 2000 {
				L = 255
			}
			goType, tsType = "string", "string"
			gormType = fmt.Sprintf("varchar(%d)", L)
		case "text":
			goType, tsType = "string", "string"
			gormType = "text"
		case "int":
			goType, tsType = "int", "number"
			gormType = "int"
		case "bigint":
			goType, tsType = "int64", "number"
			gormType = "bigint"
		case "tinyint":
			goType, tsType = "int", "number"
			gormType = "tinyint"
		case "datetime":
			goType, tsType = "time_util.LocalTime", "string"
			gormType = "datetime(3)"
		case "date":
			goType, tsType = "time_util.LocalTime", "string"
			gormType = "date"
		case "decimal":
			goType, tsType = "float64", "number"
			sc := f.DecimalScale
			if sc <= 0 || sc > 8 {
				sc = 2
			}
			gormType = fmt.Sprintf("decimal(18,%d)", sc)
		default:
			return nil, fmt.Errorf("未实现类型: %s", dt)
		}
		null := "not null"
		if f.Nullable {
			null = "null"
		}
		gormCol := fmt.Sprintf("column:%s;type:%s;%s", f.ColumnName, gormType, null)
		cm := strings.ReplaceAll(f.Comment, "`", "'")
		if cm != "" {
			gormCol += ";comment:" + cm
		}
		bind := ""
		if !f.Nullable && dt != "text" {
			bind = ` binding:"required"`
		}
		structTags := fmt.Sprintf(`gorm:"%s" json:"%s"%s`, gormCol, gj, bind)

		out = append(out, GenField{
			ColumnName: f.ColumnName,
			JSONName:   gj,
			GoName:     gn,
			GoType:     goType,
			StructTags: structTags,
			TSType:     tsType,
			ListShow:   f.ListShow,
			FormShow:   f.FormShow,
			IsQuery:    f.IsQuery,
			QueryType:  strings.ToLower(f.QueryType),
			DbType:     dt,
			Comment:    f.Comment,
		})
	}
	return out, nil
}

func filterGoType(f GenField) string {
	switch f.DbType {
	case "varchar", "text", "datetime", "date":
		return "string"
	case "int", "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "decimal":
		return "float64"
	default:
		return "string"
	}
}

func queryBlocks(_ string, qf []GenField) []string {
	var blocks []string
	for _, f := range qf {
		ft := filterGoType(f)
		col := f.ColumnName
		var b strings.Builder
		switch ft {
		case "string":
			if f.QueryType == "like" {
				b.WriteString(fmt.Sprintf("\tif req.%s != \"\" {\n\t\tdb = db.Where(\"`%s` LIKE ?\", \"%%\"+req.%s+\"%%\")\n\t}\n", f.GoName, col, f.GoName))
			} else {
				b.WriteString(fmt.Sprintf("\tif req.%s != \"\" {\n\t\tdb = db.Where(\"`%s` = ?\", req.%s)\n\t}\n", f.GoName, col, f.GoName))
			}
		case "int":
			b.WriteString(fmt.Sprintf("\tif req.%s != 0 {\n\t\tdb = db.Where(\"`%s` = ?\", req.%s)\n\t}\n", f.GoName, col, f.GoName))
		case "int64":
			b.WriteString(fmt.Sprintf("\tif req.%s != 0 {\n\t\tdb = db.Where(\"`%s` = ?\", req.%s)\n\t}\n", f.GoName, col, f.GoName))
		case "float64":
			b.WriteString(fmt.Sprintf("\tif req.%s != 0 {\n\t\tdb = db.Where(\"`%s` = ?\", req.%s)\n\t}\n", f.GoName, col, f.GoName))
		default:
			b.WriteString(fmt.Sprintf("\tif req.%s != \"\" {\n\t\tdb = db.Where(\"`%s` = ?\", req.%s)\n\t}\n", f.GoName, col, f.GoName))
		}
		blocks = append(blocks, b.String())
	}
	return blocks
}

func buildTemplateData(form *SysOnlineForm, rawFields []SysOnlineFormField) (*genTemplateData, error) {
	fields, err := buildGenFields(form, rawFields)
	if err != nil {
		return nil, err
	}
	var qf []GenField
	for _, f := range fields {
		if f.IsQuery {
			f2 := f
			f2.GoType = filterGoType(f)
			qf = append(qf, f2)
		}
	}
	hasTime := false
	for _, f := range fields {
		if strings.Contains(f.GoType, "LocalTime") {
			hasTime = true
			break
		}
	}
	qb := queryBlocks(form.EntityName, qf)

	var qts []tsFieldVue
	for i := range qf {
		f := qf[i]
		lab := f.Comment
		if lab == "" {
			lab = f.ColumnName
		}
		ft := filterGoType(f)
		ts := "string"
		if ft == "int" || ft == "int64" || ft == "float64" {
			ts = "number"
		}
		qts = append(qts, tsFieldVue{JSONName: f.JSONName, TSType: ts, Comment: lab})
	}

	var listCols []listColVue
	for _, f := range fields {
		if !f.ListShow {
			continue
		}
		t := f.Comment
		if t == "" {
			t = f.ColumnName
		}
		listCols = append(listCols, listColVue{Title: t, JSONName: f.JSONName, Ellipsis: f.DbType == "text"})
	}

	var formVue []formFieldVue
	for _, f := range fields {
		if !f.FormShow {
			continue
		}
		lab := f.Comment
		if lab == "" {
			lab = f.ColumnName
		}
		var ctrl string
		switch f.DbType {
		case "text":
			ctrl = fmt.Sprintf(`<NInput v-model:value="formModel.%s" type="textarea" placeholder="%s" :rows="3" />`, f.JSONName, lab)
		case "int", "bigint", "tinyint", "decimal":
			ctrl = fmt.Sprintf(`<NInputNumber v-model:value="formModel.%s" style="width: 100%%" placeholder="%s" />`, f.JSONName, lab)
		default:
			ctrl = fmt.Sprintf(`<NInput v-model:value="formModel.%s" placeholder="%s" />`, f.JSONName, lab)
		}
		formVue = append(formVue, formFieldVue{Label: lab, ControlTpl: ctrl})
	}

	viewPath := filepath.ToSlash(filepath.Join("generated", form.PhysTableName))
	listTitle := strings.TrimSpace(form.Description)
	if listTitle == "" {
		listTitle = form.EntityName
	}
	listTitle = strings.ReplaceAll(listTitle, `"`, "'")
	return &genTemplateData{
		PackageName:   form.PhysTableName,
		PackageImport: form.PhysTableName,
		TableName:     form.PhysTableName,
		EntityName:    form.EntityName,
		RouteGroup:    form.RouteGroup,
		Description:   form.Description,
		ServiceVar:    serviceVar(form.EntityName),
		HasTimeType:   hasTime,
		Fields:        fields,
		QueryFields:   qf,
		QueryBlocks:   qb,
		QueryFieldsTS: qts,
		ListColumns:   listCols,
		FormFieldsVue: formVue,
		ViewPath:      viewPath,
		ListTitle:     listTitle,
		ShowSearch:    len(qts) > 0,
	}, nil
}

func execTpl(name string, data any) ([]byte, error) {
	b, err := tmplFS.ReadFile(filepath.Join("templates", name))
	if err != nil {
		return nil, err
	}
	t, err := template.New(name).Delims("[[", "]]").Parse(string(b))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ZipGeneratedCode 生成 ZIP 字节（backend + frontend + 说明）
func ZipGeneratedCode(form *SysOnlineForm, fields []SysOnlineFormField) ([]byte, error) {
	data, err := buildTemplateData(form, fields)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	write := func(path string, content []byte) error {
		w, err := zw.Create(path)
		if err != nil {
			return err
		}
		_, err = w.Write(content)
		return err
	}

	modelB, err := execTpl("model.go.tpl", data)
	if err != nil {
		return nil, err
	}
	svcB, err := execTpl("service.go.tpl", data)
	if err != nil {
		return nil, err
	}
	rtB, err := execTpl("router.go.tpl", data)
	if err != nil {
		return nil, err
	}
	apiB, err := execTpl("index.api.ts.tpl", data)
	if err != nil {
		return nil, err
	}
	vueB, err := execTpl("index.vue.tpl", data)
	if err != nil {
		return nil, err
	}
	regB, err := execTpl("REGISTER.txt.tpl", data)
	if err != nil {
		return nil, err
	}

	pkg := data.PackageName
	base := "backend/api/" + pkg + "/"
	if err := write(base+"model.go", modelB); err != nil {
		return nil, err
	}
	if err := write(base+"service.go", svcB); err != nil {
		return nil, err
	}
	if err := write(base+"router.go", rtB); err != nil {
		return nil, err
	}
	if err := write("backend/REGISTER.txt", regB); err != nil {
		return nil, err
	}

	fe := "frontend/src/views/" + data.ViewPath + "/"
	if err := write(fe+"index.api.ts", apiB); err != nil {
		return nil, err
	}
	if err := write(fe+"index.vue", vueB); err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
