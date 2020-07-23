package main

import (
	"github.com/go-xorm/xorm"
	"github.com/handsomestWei/go-annotation/annotation/transaction"
	"github.com/handsomestWei/go-annotation/example"
)

// go build -gcflags=-l main.go
// main
func main() {
	scanPath := `F:\GOPATH\src\github.com\handsomestWei\go-annotation\example`
	transaction.NewTransactionManager(transaction.TransactionConfig{ScanPath: scanPath}).RegisterDao(new(example.ExampleDao))

	dao := new(example.ExampleDao)
	dao.Select()
	dao.Update(new(xorm.Session), "") // auto commit
	dao.Delete(new(xorm.Session)) // handle fail and auto rollback
}
