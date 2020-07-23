package transaction

import (
	"github.com/go-xorm/xorm"
	"github.com/handsomestWei/go-annotation/aop"
	"reflect"
)

var methodSessionMap = make(map[string]*joinPointSessionInfo)

type Transactional struct {
	ReadOnly                bool
	RollbackFor             reflect.Type
	RollbackForStructName   reflect.Type
	NoRollbackFor           reflect.Type
	NoRollbackForStructName reflect.Type
	Propagation             Propagation
	Isolation               int
	Timeout                 int
}

func (t *Transactional) Before(point *aop.JoinPoint, methodLocation string) bool {
	if methodSessionMap[methodLocation] != nil {
		t.doSessionBegin(point.Params[methodSessionMap[methodLocation].ParamSessionPosition].Interface())
	} else {
		// cache
		for i, v := range point.Params {
			if t.doSessionBegin(v.Interface()) {
				methodSessionMap[methodLocation] = &joinPointSessionInfo{
					ParamSessionPosition: i,
				}
				break
			}
		}
	}
	return true
}

func (t *Transactional) After(point *aop.JoinPoint, methodLocation string) {
	if methodSessionMap[methodLocation] != nil {
		// TODO 规约：返回值第一个参数为处理结果，类型为布尔型。由此确认是提交还是回滚
		for i, v := range point.Result {
			if i == 0 {
				switch result := v.Interface().(type) {
				case bool:
					if result {
						t.doSessionCommit(point.Params[methodSessionMap[methodLocation].ParamSessionPosition])
					} else {
						t.doSessionRollback(point.Params[methodSessionMap[methodLocation].ParamSessionPosition])
					}
				}
			} else {
				break
			}
		}
	}
}

func (t *Transactional) Finally(point *aop.JoinPoint, methodLocation string) {
	if methodSessionMap[methodLocation] != nil {
		t.doSessionClose(point.Params[methodSessionMap[methodLocation].ParamSessionPosition])
	} else {

	}
}

func (t *Transactional) IsMatch(methodLocation string) bool {
	if methodLocationMap[methodLocation] != nil {
		return true
	} else {
		return false
	}
}

func (t *Transactional) doSessionBegin(v interface{}) bool {
	switch ses := v.(type) {
	case *xorm.Session:
		ses.Begin()
		return true
	default:
		return false
	}
}

func (t *Transactional) doSessionCommit(v interface{}) bool {
	switch ses := v.(type) {
	case *xorm.Session:
		ses.Commit()
		return true
	default:
		return false
	}
}

func (t *Transactional) doSessionRollback(v interface{}) bool {
	switch ses := v.(type) {
	case *xorm.Session:
		ses.Rollback()
		return true
	default:
		return false
	}
}

func (t *Transactional) doSessionClose(v interface{}) bool {
	switch ses := v.(type) {
	case *xorm.Session:
		ses.Close()
		return true
	default:
		return false
	}
}
