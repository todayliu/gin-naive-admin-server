package config

type CaptchaConfig struct {
	KeyLong         int    `mapstructure:"key-long" json:"key-long" yaml:"key-long"`       // 验证码长度
	ImgWidth        int    `mapstructure:"img-width" json:"img-width" yaml:"img-width"`    // 验证码宽度
	ImgHeight       int    `mapstructure:"img-height" json:"img-height" yaml:"img-height"` // 验证码高度
	NoiseCount      int    `mapstructure:"noise-count" json:"noise-count" yaml:"noise-count"`
	ShowLineOptions int    `mapstructure:"show-line-options" json:"show-line-options" yaml:"show-line-options"`
	Source          string `mapstructure:"source" json:"source" yaml:"source"`
}
