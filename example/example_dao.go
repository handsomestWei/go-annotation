package example

import "github.com/go-xorm/xorm"

type ExampleDao struct {
}

//select
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
