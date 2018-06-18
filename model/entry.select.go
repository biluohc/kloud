package model

import (
	"time"

	log "github.com/cihub/seelog"
)

func FindEntryByID(id string, leftJoinFile bool) (*Entry, error) {
	var entry = new(Entry)
	var err error
	if leftJoinFile {
		err = Db.Model(entry).
			Relation("File").
			Where("entry.id = ?", id).
			Select()
	} else {
		err = Db.Model(entry).
			Where("id = ?", id).
			Select()
	}
	if err != nil {
		return nil, err
	}
	return entry, err
}

func FindEntryByIDandUid(id string, uid string, leftJoinFile bool) (*Entry, error) {
	var entry = new(Entry)
	var err error
	if leftJoinFile {
		err = Db.Model(entry).
			Relation("File").
			Where("entry.id = ? and entry.user_id = ?", id, uid).
			Select()
	} else {
		err = Db.Model(entry).
			Where("id = ? and user_id = ?", id, uid).
			Select()
	}
	if err != nil {
		return nil, err
	}
	return entry, err
}

// 用户拥有的同一文件数量
func CountEntryByFidAndUid(fid string, uid string) (int, error) {
	return Db.Model((*Entry)(nil)).
		Where("user_id =?", uid).
		Count()
}

func FindOrCreateRootByUid(u *User) (*Entry, error) {
	var entry = new(Entry)
	err := Db.Model(entry).
		Where("user_id = ? and name='/' and parent_id is null", u.ID).
		Select()

	if err != nil {
		new := Entry{
			UserID:     u.ID,
			FileID:     "",
			ParentID:   "",
			Name:       "/",
			CreateTime: time.Now(),
			ModifyTime: time.Now(),
		}
		err = Db.Insert(&new)
		return &new, err
	}

	return entry, err
}

func (e *Entry) Parent() (*Entry, error) {
	var entry = new(Entry)
	err := Db.Model(entry).
		Where("id = ?", e.ParentID).
		Select()
	return entry, err
}

func FindEntryParentByIDandUid(id string, uid string) (*Entry, error) {
	var entry = new(Entry)

	// 这个正确性有待验证。。
	tmpEntry := Db.Model((*Entry)(nil)).Where("id = ? and user_id = ?", id, uid)
	err := Db.Model((*Entry)(nil)).TableExpr("(?)", tmpEntry).Select(entry)
	log.Info(entry, err)
	return entry, err
}

func FindEntryPath(id string) ([]Entry, error) {
	var es []Entry
	_, err := Db.Query(&es, `with RECURSIVE EntryTree AS
		(
			SELECT * from kloud_entry where Id=? 
			UNION ALL
			SELECT kloud_entry.* from EntryTree
			JOIN kloud_entry on EntryTree.parent_id= kloud_entry.id
		)
		SELECT * from EntryTree;`, id)
	return es, err
}

func (e *Entry) Subs(leftJoinFile bool) ([]Entry, error) {
	var es []Entry
	var err error
	if leftJoinFile {
		err = Db.Model(&es).
			// trash 判断
			Where("parent_id = ?  and delete_time is null", e.ID).
			Column("entry.*", "File").
			Select()
	} else {
		err = Db.Model(&es).
			Where("parent_id = ?  and delete_time is null", e.ID).
			Select()
	}
	return es, err
}

func FindEntrySubsByIDandUid(id string, uid string, leftJoinFile bool) ([]Entry, error) {
	var es []Entry
	var err error
	if leftJoinFile {
		err = Db.Model(&es).
			// trash 判断
			Where("parent_id = ? and user_id = ? and delete_time is null", id, uid).
			Column("entry.*", "File").
			Select()
	} else {
		err = Db.Model(&es).
			// trash 判断
			Where("parent_id = ? and user_id = ? and delete_time is null", id, uid).
			Select()
	}
	return es, err
}

// 这个orm好像不能这么搞， 等换了再说。。

// -- mysql
// with EntryTree AS
// (
//     SELECT * from kloud_entry where Id=? --需要查找的子节点
//     UNION ALL
//     SELECT kloud_entry.* from EntryTree
//     JOIN kloud_entry on EntryTree.parent_id= kloud_entry.id
// )
// SELECT * from EntryTree;

// -- pg
