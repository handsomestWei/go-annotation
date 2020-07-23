package transaction

import (
	"fmt"
	"github.com/handsomestWei/go-annotation/aop"
	"github.com/handsomestWei/go-annotation/ast"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"unicode"
)

var config TransactionConfig
var methodLocationMap = make(map[string]*struct{})

type joinPointSessionInfo struct {
	ParamSessionPosition int
}

type TransactionConfig struct {
	ScanPath string
}

type transactionManager struct {
}

func NewTransactionManager(cfg TransactionConfig) *transactionManager {
	if err := cfg.check(); err != nil {
		panic(err)
	}
	config = cfg
	scanGoFile()
	aop.RegisterAspect(new(Transactional))
	return new(transactionManager)
}

func (c *TransactionConfig) check() error {
	_, err := os.Stat(c.ScanPath)
	return err
}

// TODO
func (c *TransactionConfig) Reload() error {
	return nil
}

func (t *transactionManager) RegisterDao(daoLs ...interface{}) (tm *transactionManager) {
	tm = t
	for _, v := range daoLs {
		aop.RegisterPoint(reflect.TypeOf(v))
	}
	return
}

func scanGoFile() {
	filepath.Walk(config.ScanPath, walkFunc)
}

func walkFunc(fullPath string, info os.FileInfo, err error) error {
	if info == nil {
		return err
	}
	if info.IsDir() {
		return nil
	} else {
		if GO_FILE_SUFIX == path.Ext(fullPath) {
			// is go file
			cacheMethodLocationMap(fullPath)
		}
		return nil
	}
}

func cacheMethodLocationMap(fullPath string) {
	bt, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return
	}
	result := ast.ScanFuncDeclByComment("", string(bt), COMMENT_NAME)
	if result == nil {
		return
	}
	if result.RecvMethods != nil {
		for _, l := range result.RecvMethods {
			for _, v := range l {
				for i, s := range v.MethodName {
					// skip the private method
					if i == 0 && unicode.IsUpper(s) {
						methodLocationMap[fmt.Sprintf("%s.%s.%s", v.PkgName, v.RecvName, v.MethodName)] = new(struct{})
					}
					break
				}
			}
		}
	}
}
