package captcha

import (
	"image/color"

	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

// GenerateCaptchaHandler base64Captcha create http handler
func GenerateCaptchaHandler() (id, b64s string, err error) {
	driverString := new(base64Captcha.DriverString)
	driverString.Length = 6
	driverString.Source = "1234567890zxcvbnmlkjhgfdsaqwertyuiop"
	driverString.Width = 220
	driverString.Height = 47
	driverString.BgColor = &color.RGBA{0, 0, 0, 0}
	driverString.Fonts = []string{"wqy-microhei.ttc"}
	c := base64Captcha.NewCaptcha(driverString.ConvertFonts(), store)
	return c.Generate()
}

// CaptchaVerifyHandle base64Captcha verify http handler
func CaptchaVerifyHandle(id, val string) bool {
	return store.Verify(id, val, true)
}
