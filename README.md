# go-annotation
golang 使用注释实现类似java的注解机制。基于ast语法解析和monkey动态代理。目前实现`@Transactional`的demo

# Usage

在DAO层使用`//@Transactional`注释标记目标方法。自动实现事务处理，不用额外编写事务处理代码。

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
目前无法实现类似java的Class.ForName("class full name")字符串转对象，只能显式调用RegisterDao方法传入  
编译时需禁用内联go build -gcflags=-l  

```
// go build -gcflags=-l main.go
// main
func main() {
	scanPath := `xxx\github.com\handsomestWei\go-annotation\example`
	// 初始化事务管理器，扫描指定包路径，代理DAO对象
	transaction.NewTransactionManager(transaction.TransactionConfig{ScanPath: scanPath}).RegisterDao(new(example.ExampleDao))

	dao := new(example.ExampleDao)
	dao.Select()
	// 事务自动开启。处理成功将自动提交和关闭事务，处理失败将自动回滚事务
	dao.Update(new(xorm.Session), "")
	dao.Delete(new(xorm.Session))
}
```
