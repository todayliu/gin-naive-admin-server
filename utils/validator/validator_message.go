package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func GetValidatorErrorMessage(err error, obj interface{}) string {
	// 将 err 转型为 validator.ValidationErrors
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			// 通过反射获取该结构体的类型
			t := reflect.TypeOf(obj)
			// 如果是指针，获取其指向的元素
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}

			// 找到报错的字段
			if field, ok := t.FieldByName(e.Field()); ok {
				// 读取该字段上的 msg 标签内容
				message := field.Tag.Get("message")
				if message != "" {
					return message
				}
			}
		}
	}
	return err.Error() // 如果没找到自定义 msg，返回默认错误
}
