package config

// SecurityConfig API 权限与其它安全相关开关
type SecurityConfig struct {
	// SuperRoleCodes 额外超级角色编码（与角色 admin 合并）；拥有任一者可视为超管（如全部按钮权限）
	SuperRoleCodes []string `mapstructure:"super-role-codes" json:"super-role-codes" yaml:"super-role-codes"`
	// RelaxApiAuth 为 true 时仅校验登录，不校验权限码（便于未初始化按钮权限的库）
	RelaxApiAuth bool `mapstructure:"relax-api-auth" json:"relax-api-auth" yaml:"relax-api-auth"`
}
