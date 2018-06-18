package model

import (
	"errors"
	"time"

	log "github.com/cihub/seelog"
	"github.com/go-pg/pg/orm"
)

type Public struct {
	tableName  struct{}  `sql:"kloud_public"`
	ID         string    `sql:"id,pk" json:"id"`
	Name       string    `sql:"name, notnull" json:"name"`
	UserID     string    `sql:"user_id,notnull" json:"userId"`
	EntryID    string    `sql:"entry_id" json:"entryId"`
	CreateTime time.Time `db:"create_time,notnull" json:"createTime"`
	ExpireTime time.Time `db:"expire_time" json:"expireTime"`
	Referer    string    `db:"referer" json:"referer"`
}

func (p *Public) BeforeInsert(db orm.DB) error {
	if p.ID == "" {
		p.ID = NanoId()
	}
	return nil
}

func (u *User) Publics() ([]Public, error) {
	var ps = make([]Public, 0)
	err := Db.Model(&ps).
		Where("user_id = ?", u.ID).
		Select()
	return ps, err
}

func FindPublicByID(id string) (*Public, error) {
	var p = new(Public)
	err := Db.Model(p).
		Where("id = ?", id).
		Select()

	var zeroTime time.Time
	if zeroTime != p.ExpireTime && p.ExpireTime.Before(time.Now()) {
		err = errors.New("Session过期")
		_ = p.Delete()
	}
	return p, err
}

func FindPublicByIDandUid(id string, uid string) (*Public, error) {
	var s = new(Public)
	err := Db.Model(s).
		Where("id = ? and user_id =? ", id, uid).
		Select()

	var zeroTime time.Time
	if zeroTime != s.ExpireTime && s.ExpireTime.Before(time.Now()) {
		err = errors.New("Session过期")
		_ = s.Delete()
	}

	return s, err
}

func (e *Entry) Public() (*Public, error) {
	new := Public{
		UserID:     e.UserID,
		EntryID:    e.ID,
		CreateTime: time.Now(),
		Name:       e.Name,
	}
	err := Db.Insert(&new)
	return &new, err
}

func (p *Public) Delete() error {
	_, err := Db.Model(p).Delete()
	return err
}
func (s *Public) SetReferer(referer string) error {
	res, err := Db.Model(s).
		Set("referer = ?", referer).
		Where("id = ?", s.ID).
		Update()
	if err != nil {
		log.Errorf("SetReferer for public(%v)  affectedRows(%d) failed: %v", s, res.RowsAffected(), err)
	}
	return err
}

func (s *Public) SetMaxAge(maxAge int) error {
	var etime time.Time
	if maxAge > 0 {
		etime = time.Now().Add(time.Duration(maxAge) * time.Second)
	}

	res, err := Db.Model(s).
		Set("expire_time = ?", etime).
		Where("id = ?", s.ID).
		Update()
	if err != nil {
		log.Errorf("SetMaxAge for public(%v)  affectedRows(%d) failed: %v", s, res.RowsAffected(), err)
	}
	return err
}

func (e *Entry) PublicWithReferers(referer string, maxAge int) (*Public, error) {
	new := Public{
		UserID:     e.UserID,
		EntryID:    e.ID,
		CreateTime: time.Now(),
		Referer:    referer,
		Name:       e.Name,
	}
	if maxAge > 0 {
		new.ExpireTime = time.Now().Add(time.Duration(maxAge) * time.Second)
	}
	err := Db.Insert(&new)
	return &new, err
}
