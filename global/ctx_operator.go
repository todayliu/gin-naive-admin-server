package global

import "context"

type operatorUIDKey struct{}

var operatorUIDCtxKey = operatorUIDKey{}

// WithOperatorUserID 将当前操作人用户 ID 写入 context（供 GORM 审计字段使用）
func WithOperatorUserID(ctx context.Context, uid uint) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, operatorUIDCtxKey, uid)
}

// OperatorUserID 从 context 读取操作人用户 ID；未设置或公开接口返回 0
func OperatorUserID(ctx context.Context) uint {
	if ctx == nil {
		return 0
	}
	v := ctx.Value(operatorUIDCtxKey)
	if v == nil {
		return 0
	}
	if id, ok := v.(uint); ok {
		return id
	}
	return 0
}
