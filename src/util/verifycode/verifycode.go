package verifycode

import (
	cap "github.com/mojocn/base64Captcha"
)

const (
	Cap_Num_Mod     = cap.CaptchaModeNumber
	Cap_Char_Mod    = cap.CaptchaModeAlphabet
	Cap_Metic_Mod   = cap.CaptchaModeArithmetic
	Cap_NumChar_Mod = cap.CaptchaModeNumberAlphabet
)

/*
	60x240
	mode means verifycode type
	0:Number, 1:Alphabet, 2:Arithmetic, 3:NumberAlphabet
*/
func CodeGenerate(heigt, width, mode int) (capId, pngBase64 string) {

	configChar := cap.ConfigCharacter{
		Height:             heigt,
		Width:              width,
		Mode:               mode,
		ComplexOfNoiseText: cap.CaptchaComplexLower,
		ComplexOfNoiseDot:  cap.CaptchaComplexLower,
		IsUseSimpleFont:    false,
		IsShowHollowLine:   true,
		IsShowNoiseDot:     true,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         4,
	}

	capId, charCap := cap.GenerateCaptcha("", configChar)

	pngBase64 = cap.CaptchaWriteToBase64Encoding(charCap)

	return
}

func CodeGenerateByCapId(heigt, width, mode int, capId string) (pngBase64 string) {

	configChar := cap.ConfigCharacter{
		Height:             heigt,
		Width:              width,
		Mode:               mode,
		ComplexOfNoiseText: cap.CaptchaComplexLower,
		ComplexOfNoiseDot:  cap.CaptchaComplexLower,
		IsUseSimpleFont:    false,
		IsShowHollowLine:   true,
		IsShowNoiseDot:     true,
		IsShowNoiseText:    true,
		IsShowSlimeLine:    true,
		IsShowSineLine:     true,
		CaptchaLen:         4,
	}

	_, charCap := cap.GenerateCaptcha(capId, configChar)

	pngBase64 = cap.CaptchaWriteToBase64Encoding(charCap)

	return
}

func CodeValidate(identifier, verifyValue string) bool {
	return cap.VerifyCaptcha(identifier, verifyValue)
}
