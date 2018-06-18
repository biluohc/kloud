package model

import (
	"time"

	"github.com/go-pg/pg/orm"
)

// trash 判断.
func FindTrashesByUid(uid string) ([]Entry, error) {
	var es []Entry
	err := Db.Model(&es).
		Where("user_id = ?  and delete_time is not null", uid).
		Column("entry.*", "File").
		Select()
	return es, err
}

func (e *Entry) Trash() (*Entry, orm.Result, error) {
	res, err := Db.Model(e).
		Set("delete_time = ?", time.Now()).
		Where("id = ?", e.ID).
		Update()
	return e, res, err
}

func TrashbyIDandUid(id string, uid string) (*Entry, orm.Result, error) {
	e := new(Entry)
	res, err := Db.Model(e).
		Set("delete_time = ?", time.Now()).
		Where("id = ? and user_id= ?", id, uid).
		Update()
	return e, res, err
}

func (e *Entry) UnTrash() (*Entry, orm.Result, error) {
	res, err := Db.Model(e).
		// Set("delete_time = ?", 0). // 0 being become null by gopg
		Set("delete_time = null").
		Where("id = ?", e.ID).
		Update()
	return e, res, err
}

func UnTrashbyIDandUid(id string, uid string) (*Entry, orm.Result, error) {
	e := new(Entry)
	res, err := Db.Model(e).
		// Set("delete_time = ?", 0). // 0 being become null by gopg
		Set("delete_time = null").
		Where("id = ? and user_id = ?", id, uid).
		Update()
	return e, res, err
}
