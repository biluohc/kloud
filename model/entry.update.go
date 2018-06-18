package model

import (
	"github.com/go-pg/pg/orm"
)

// newPid下有同名文件或者会有人把自己的挂到别人的 entry 里， 由调用者检查
func MoveEntryByIDandUid(newName string, newPid string, eid string, uid string) (*Entry, orm.Result, error) {
	e := new(Entry)
	if newName == "" {
		res, err := Db.Model(e).
			Set("parent_id = ?", newPid).
			Where("id = ? AND user_id =?", eid, uid).
			Update()
		return e, res, err
	}
	res, err := Db.Model(e).
		Set("parent_id = ?", newPid).
		Set("name = ?", newName).
		Where("id = ? AND user_id =?", eid, uid).
		Update()
	return e, res, err
}

// newPid下有同名文件， 由调用者检查
func RenameEntryByIDandUid(newName string, eid string, uid string) (*Entry, orm.Result, error) {
	e := new(Entry)
	res, err := Db.Model(e).
		Set("name = ?", newName).
		Where("id = ? AND user_id =?", eid, uid).
		Update()
	return e, res, err
}
