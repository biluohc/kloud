package model

import (
	"errors"
	"time"

	log "github.com/cihub/seelog"
	"github.com/go-pg/pg/orm"
)

type Share struct {
	tableName  struct{}  `sql:"kloud_share"`
	ID         string    `sql:"id,pk" json:"id"`
	Name       string    `sql:"name, notnull" json:"name"`
	UserID     string    `sql:"user_id,notnull" json:"userId"`
	EntryID    string    `sql:"entry_id" json:"entryId"`
	CreateTime time.Time `db:"create_time,notnull" json:"createTime"`
	ExpireTime time.Time `db:"expire_time" json:"expireTime"`
	Password   string    `db:"password" json:"password"`
	Entry      *Entry    `db:"-" json:"entry"`
}

func (s *Share) BeforeInsert(db orm.DB) error {
	if s.ID == "" {
		s.ID = NanoId()
	}
	return nil
}

func (u *User) Shares() ([]Share, error) {
	// 这个的json不是 null, 没有初始化的是 null
	var ss = make([]Share, 0)
	err := Db.Model(&ss).
		Where("user_id = ?", u.ID).
		Select()
	if err != nil {
		return nil, err
	}
	return ss, err
}

func FindShareByID(id string) (*Share, error) {
	var s = new(Share)
	err := Db.Model(s).
		Where("id = ?", id).
		Select()
	if err != nil {
		return nil, err
	}

	var zeroTime time.Time
	if zeroTime != s.ExpireTime && s.ExpireTime.Before(time.Now()) {
		err = errors.New("Session过期")
		_ = s.Delete()
	}

	return s, err
}

func FindShareByIDandUid(id string, uid string) (*Share, error) {
	var s = new(Share)
	err := Db.Model(s).
		Where("id = ? and user_id =? ", id, uid).
		Select()
	if err != nil {
		return nil, err
	}

	var zeroTime time.Time
	if zeroTime != s.ExpireTime && s.ExpireTime.Before(time.Now()) {
		err = errors.New("Session过期")
		_ = s.Delete()
	}
	return s, err
}

func (e *Entry) Share() (*Share, error) {
	new := Share{
		UserID:     e.UserID,
		EntryID:    e.ID,
		CreateTime: time.Now(),
		Name:       e.Name,
	}
	err := Db.Insert(&new)
	return &new, err
}

func (s *Share) Delete() error {
	_, err := Db.Model(s).Delete()
	return err
}

func (s *Share) SetPassword(word string) error {
	res, err := Db.Model(s).
		Set("password = ?", word).
		Where("id = ?", s.ID).
		Update()
	if err != nil {
		log.Errorf("SetPassword for Share(%v)  affectedRows(%d) failed: %v", s, res.RowsAffected(), err)
	}
	return err
}

func (s *Share) SetMaxAge(maxAge int) error {
	var etime time.Time
	if maxAge > 0 {
		etime = time.Now().Add(time.Duration(maxAge) * time.Second)
	}

	res, err := Db.Model(s).
		Set("expire_time = ?", etime).
		Where("id = ?", s.ID).
		Update()
	if err != nil {
		log.Errorf("SetMaxAge for Share(%v)  affectedRows(%d) failed: %v", s, res.RowsAffected(), err)
	}
	return err
}

func (e *Entry) ShareWithPassWordAndMaxAge(pw string, maxAge int) (*Share, error) {
	new := Share{
		UserID:     e.UserID,
		EntryID:    e.ID,
		CreateTime: time.Now(),
		Password:   pw,
		Name:       e.Name,
	}
	if maxAge > 0 {
		new.ExpireTime = time.Now().Add(time.Duration(maxAge) * time.Second)
	}

	err := Db.Insert(&new)
	return &new, err
}
