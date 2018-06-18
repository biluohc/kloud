package model

import (
	"time"

	log "github.com/cihub/seelog"
	"github.com/go-pg/pg/orm"
)

type Sess struct {
	tableName  struct{}  `sql:"kloud_session"`
	ID         string    `sql:"id,pk" json:"id"`
	UserID     string    `sql:"user_id" json:"userId"`
	UserAgent  string    `sql:"user_agent, notnull" json:"userAgent"`
	IP         string    `sql:"ip, notnull" json:"ip"`
	CreateTime time.Time `db:"create_time" json:"createTime"`
	ExpireTime time.Time `db:"expire_time" json:"expireTime"`
}

func (u *Sess) BeforeInsert(db orm.DB) error {
	if u.ID == "" {
		u.ID = NanoId()
	}
	return nil
}

func FindSessByID(id string) (*Sess, error) {
	var u = new(Sess)
	err := Db.Model(u).
		Where("id = ?", id).
		Select()
	// if err == nil && u.ExpireTime.Before(time.Now()) {
	// 	err = errors.New("Session过期")
	// 	_ = u.Delete()
	// }
	return u, err
}

func DeleteSessByID(id string) error {
	var se = new(Sess)
	res, err := Db.Model(se).
		Where("id = ?", id).
		Delete()
	if res.RowsAffected() > 1 {
		log.Error("Sess.DeleteByID-RowsAffected()>1: %d %v", res.RowsAffected(), err)
	}
	return err
}

func (sess *Sess) Delete() error {
	_, err := Db.Model(sess).Delete()
	if err != nil {
		log.Error("Sess.Delete", err)
	}
	return err
}

func (sess *Sess) User() (*User, error) {
	return FindUserByID(sess.UserID)
}

// 86400 = 1 day
// MaxAge=0 means no Max-Age attribute specified and the cookie will be
// deleted after the browser session ends.
// MaxAge<0 means delete cookie immediately.
// MaxAge>0 means Max-Age attribute present and given in seconds.
func CreateSess(uid string, maxAge int, ua string, ip string) (*Sess, error) {
	var s = Sess{
		UserID:     uid,
		UserAgent:  ua,
		IP:         ip,
		CreateTime: time.Now(),
		ExpireTime: time.Now().Add(time.Duration(maxAge) * time.Second),
	}

	err := Db.Insert(&s)
	return &s, err
}
