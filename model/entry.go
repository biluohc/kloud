package model

import (
	"time"

	"github.com/go-pg/pg/orm"
)

type Entry struct {
	tableName  struct{}  `sql:"kloud_entry"`
	ID         string    `sql:"id,pk" json:"id"`
	UserID     string    `sql:"user_id,notnull" json:"userId"`
	FileID     string    `sql:"file_id" json:"fileId"`
	File       *File     `sql:"-" json:"file"`
	ParentID   string    `sql:"parent_id" json:"parentId"`
	Name       string    `sql:"name, notnull" json:"name"`
	CreateTime time.Time `db:"create_time,notnull" json:"createTime"`
	ModifyTime time.Time `db:"modify_time,notnull" json:"modifyTime"`
	DeleteTime time.Time `db:"delete_time" json:"deleteTime"`
}

func (entry *Entry) BeforeInsert(db orm.DB) error {
	if entry.ID == "" {
		entry.ID = NanoId()
	}
	return nil
}

// '' fid is null
func CreateEntryByUid(name string, pid string, fid string, uid string) (i *Entry, err error) {
	// 调用者要检查非法名字
	// input := []byte(name)
	// encodeString := base64.StdEncoding.EncodeToString(input)
	new := Entry{
		UserID:     uid,
		FileID:     fid,
		ParentID:   pid,
		Name:       name,
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
	}
	err = Db.Insert(&new)
	return &new, err
}
