# go-annotation
golang注释实现类似java的注解机制。基于ast语法解析和monkey动态代理。目前实现`@Transactional`的demo

# Usage

使用`//@Transactional`注释标记目标方法，自动实现事务处理

```
	type ExampleDao struct {
    }

	func (e *ExampleDao) Select() (bool, error) {
		return true, nil
	}

	//@Transactional
	func (d *ExampleDao) Update(s *xorm.Session, param string) (bool, error) {
		return true, nil
	}
	
	//@Transactional
	func (d *ExampleDao) Delete(s *xorm.Session) (bool, error) {
		return false, nil
	}
```

事务管理器`TransactionManager`启动，遍历go文件，获取被指定注释标记的`包名.接收者名.方法名`
目前无法实现类似java的Class.ForName("class full name")字符串转对象，只能显式调用RegisterDao方法传入对象

```
    // go build -gcflags=-l main.go
	// main
	func main() {
		scanPath := `/xxx/xxx` // scan file director
		transaction.NewTransactionManager(transaction.TransactionConfig{ScanPath: scanPath}).RegisterDao(new(example.ExampleDao))

		dao := new(example.ExampleDao)
		dao.Select()
		dao.Update(new(xorm.Session), "") // auto commit and close
		dao.Delete(new(xorm.Session)) // handle fail and auto rollback
	}
```
