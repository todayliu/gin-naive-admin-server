package global

import (
	"reflect"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// FillAuditDisplayNames 根据 createBy/updateBy（sys_user.id）批量填充 CreateByName、UpdateByName（u_name）。
// dest 支持 *T（单条）或 *[]T / *[]*T（切片）；T 需内嵌 GNA_MODEL 或包含同名导出字段。
func FillAuditDisplayNames(db *gorm.DB, dest interface{}) {
	if db == nil || dest == nil {
		return
	}
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}
	rv = rv.Elem()

	var elems []reflect.Value
	switch rv.Kind() {
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			el := rv.Index(i)
			if el.Kind() == reflect.Ptr {
				if el.IsNil() {
					continue
				}
				el = el.Elem()
			}
			if el.Kind() == reflect.Struct {
				elems = append(elems, el)
			}
		}
	case reflect.Struct:
		elems = append(elems, rv)
	default:
		return
	}

	idSet := map[uint]struct{}{}
	for _, el := range elems {
		cb, cbOk := uintFromStructField(el, "CreateBy")
		ub, ubOk := uintFromStructField(el, "UpdateBy")
		if cbOk && cb > 0 {
			idSet[cb] = struct{}{}
		}
		if ubOk && ub > 0 {
			idSet[ub] = struct{}{}
		}
	}
	if len(idSet) == 0 {
		return
	}
	ids := make([]uint, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	var rows []struct {
		ID    uint   `gorm:"column:id"`
		UName string `gorm:"column:u_name"`
	}
	if err := db.Unscoped().Table("sys_user").Select("id", "u_name").Where("id IN ?", ids).Find(&rows).Error; err != nil {
		if GNA_LOG != nil {
			GNA_LOG.Warn("FillAuditDisplayNames: 加载 u_name 失败", zap.Error(err))
		}
		return
	}
	names := make(map[uint]string, len(rows))
	for _, r := range rows {
		names[r.ID] = r.UName
	}
	for _, el := range elems {
		if cb, ok := uintFromStructField(el, "CreateBy"); ok && cb > 0 {
			setStringField(el, "CreateByName", names[cb])
		}
		if ub, ok := uintFromStructField(el, "UpdateBy"); ok && ub > 0 {
			setStringField(el, "UpdateByName", names[ub])
		}
	}
}

func uintFromStructField(v reflect.Value, name string) (uint, bool) {
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return 0, false
	}
	f := v.FieldByName(name)
	if !f.IsValid() || !f.CanInterface() {
		return 0, false
	}
	switch f.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint(f.Uint()), true
	default:
		return 0, false
	}
}

func setStringField(v reflect.Value, name, val string) {
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return
	}
	f := v.FieldByName(name)
	if !f.IsValid() || !f.CanSet() || f.Kind() != reflect.String {
		return
	}
	f.SetString(val)
}
