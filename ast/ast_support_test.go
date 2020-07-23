package ast

import (
	"fmt"
	"go/parser"
	"io/ioutil"
	"testing"
)

type Transact struct {
}

//@Transactional
func (*Transact) Before() {
}

func TestPrintAstInfo(t *testing.T) {
	bt, _ := ioutil.ReadFile(`ast_support_test.go`)
	src := string(bt)
	PrintAstInfo(``, src, parser.ParseComments)
}

func TestScanFuncDeclByComment(t *testing.T) {
	bt, _ := ioutil.ReadFile(`ast_support_test.go`)
	src := string(bt)
	result := ScanFuncDeclByComment(``, src, "@Transactional")
	fmt.Println(fmt.Sprintf("%v", result))
}
