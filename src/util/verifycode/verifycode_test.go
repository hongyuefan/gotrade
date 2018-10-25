package verifycode

import (
	"testing"
)

func TestCodeGenerate(t *testing.T) {

	id, png := CodeGenerate(60, 240, 0)

	t.Log(id, png)

	png = CodeGenerateByCapId(60, 240, 0, id)

	t.Log(png)

	ok := CodeValidate(id, "")

	t.Log(ok)
}
