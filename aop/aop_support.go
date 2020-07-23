package aop

import (
	"bou.ke/monkey"
	"fmt"
	"reflect"
	"strings"
)

type JoinPoint struct {
	Receiver interface{}
	Method   reflect.Method
	Params   []reflect.Value
	Result   []reflect.Value
}

func NewJoinPoint(receiver interface{}, params []reflect.Value, method reflect.Method) *JoinPoint {
	point := &JoinPoint{
		Receiver: receiver,
		Params:   params,
		Method:   method,
	}

	fn := method.Func
	fnType := fn.Type()
	nout := fnType.NumOut()
	point.Result = make([]reflect.Value, nout)
	for i := 0; i < nout; i++ {
		point.Result[i] = reflect.Zero(fnType.Out(i))
	}

	return point
}

type AspectInterface interface {
	Before(point *JoinPoint, methodLocation string) bool
	After(point *JoinPoint, methodLocation string)
	Finally(point *JoinPoint, methodLocation string)
	IsMatch(methodLocation string) bool
}

var aspectList = make([]AspectInterface, 0)

func RegisterPoint(pointType reflect.Type) {
	pkgPth := pointType.PkgPath()
	receiverName := pointType.Name()
	if pointType.Kind() == reflect.Ptr {
		pkgPth = pointType.Elem().PkgPath()
		receiverName = pointType.Elem().Name()
	}
	for i := 0; i < pointType.NumMethod(); i++ {
		method := pointType.Method(i)
		pkgList := strings.Split(pkgPth, "/")
		methodLocation := fmt.Sprintf("%s.%s.%s", pkgList[len(pkgList)-1], receiverName, method.Name)
		var guard *monkey.PatchGuard
		var proxy = func(in []reflect.Value) []reflect.Value {
			guard.Unpatch()
			defer guard.Restore()
			receiver := in[0]
			point := NewJoinPoint(receiver, in[1:], method)
			defer finallyProcessed(point, methodLocation)
			if !beforeProcessed(point, methodLocation) {
				return point.Result
			}
			point.Result = receiver.MethodByName(method.Name).Call(in[1:])
			afterProcessed(point, methodLocation)
			return point.Result
		}
		proxyFn := reflect.MakeFunc(method.Func.Type(), proxy)
		guard = monkey.PatchInstanceMethod(pointType, method.Name, proxyFn.Interface())
	}
}

func RegisterAspect(aspect AspectInterface) {
	aspectList = append(aspectList, aspect)
}

func beforeProcessed(point *JoinPoint, methodLocation string) bool {
	for _, aspect := range aspectList {
		if !aspect.IsMatch(methodLocation) {
			continue
		}
		if !aspect.Before(point, methodLocation) {
			return false
		}
	}
	return true
}

func afterProcessed(point *JoinPoint, methodLocation string) {
	for i := len(aspectList) - 1; i >= 0; i-- {
		aspect := aspectList[i]
		if !aspect.IsMatch(methodLocation) {
			continue
		}
		aspect.After(point, methodLocation)
	}
}

func finallyProcessed(point *JoinPoint, methodLocation string) {
	for i := len(aspectList) - 1; i >= 0; i-- {
		aspect := aspectList[i]
		if !aspect.IsMatch(methodLocation) {
			continue
		}
		aspect.Finally(point, methodLocation)
	}
}
