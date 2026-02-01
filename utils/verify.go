package utils

var (
	LoginVerify = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty(), Gt("6"), Lt("20")}, "CaptchaId": {NotEmpty(), Eq("4")}}
)
