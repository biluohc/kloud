package model

import (
	"time"

	log "github.com/cihub/seelog"
	"github.com/go-pg/pg/orm"
)

type User struct {
	tableName  struct{}  `sql:"kloud_user"`
	ID         string    `sql:"id,pk" json:"id"`
	Name       string    `db:"name,notnull" json:"name"`
	Email      string    `db:"email,notnull" json:"email"`
	Role       uint16    `db:"role,notnull" json:"role"`
	Password   string    `db:"password,notnull" json:"-"`
	AvatarURL  string    `db:"avatar_url" json:"avatarURL"`
	CreateTime time.Time `db:"create_time,notnull" json:"createTime" `
	Disabled   bool      `db:"disabled,notnull" json:"disabled"`
	MaxSize    int64     `db:"max_size,notnull" json:"maxSize"`
	CurSize    int64     `db:"cur_size,notnull" json:"curSize"`
}

func (u *User) BeforeInsert(db orm.DB) error {
	if u.ID == "" {
		u.ID = NanoId()
	}
	return nil
}

func FindUserByEmail(email string) (*User, error) {
	var u = new(User)
	err := Db.Model(u).
		Where("email = ?", email).
		Select()
	return u, err
}

func FindUserByID(uid string) (*User, error) {
	var u = new(User)
	err := Db.Model(u).
		Where("id = ?", uid).
		Select()
	return u, err
}

func GetUsers() ([]User, error) {
	var users []User
	err := Db.Model(&users).
		Select()
	return users, err
}

func (u *User) SubCurSize(size int64) error {
	res, err := Db.Model(u).
		Set("cur_size = cur_size-?", size).
		Where("id = ?", u.ID).
		Update()
	if err != nil {
		log.Errorf("SubCurSize for User(%v)  affectedRows(%d) failed: %v", u, res.RowsAffected(), err)
	}
	return err
}

// 这要么直接用 cur_size = cur_size+?
// 这要么用 cur_size = \"user\".cur_size+?
func (u *User) AddCurSize(size int64) error {
	res, err := Db.Model(u).
		Set("cur_size = cur_size+?", size).
		Where("id = ?", u.ID).
		Update()
	if err != nil {
		log.Errorf("AddCurSize for User(%v)  affectedRows(%d) failed: %v", u, res.RowsAffected(), err)
	}
	return err
}

func (u *User) Update() error {
	// 垃圾，db: notnull对 结构直接update无效。。
	// return Db.Update(u)
	_, err := Db.Model(u).
		Set("name = ?", u.Name).
		Set("email = ?", u.Email).
		Set("role = ?", u.Role).
		Set("password = ?", u.Password).
		Set("max_size = ?", u.MaxSize).
		Set("avatar_url = ?", u.AvatarURL).
		Where("id = ?", u.ID).
		Update()
	return err
}

func (u *User) SetMaxSize(size int64) error {
	res, err := Db.Model(u).
		Set("max_size = ?", size).
		Where("id = ?", u.ID).
		Update()
	if err != nil {
		log.Errorf("SetMaxSize for User(%v)  affectedRows(%d) failed: %v", u, res.RowsAffected(), err)
	}
	return err
}

func (u *User) SetName(name string) error {
	res, err := Db.Model(u).
		Set("name = ?", name).
		Where("id = ?", u.ID).
		Update()
	if err != nil {
		log.Errorf("SetName for User(%v)  affectedRows(%d) failed: %v", u, res.RowsAffected(), err)
	}
	return err
}

func (u *User) SetDisabled(status bool) error {
	res, err := Db.Model(u).
		Set("disabled = ?", status).
		Where("id = ?", u.ID).
		Update()
	if err != nil {
		log.Errorf("SetDisabled for User(%v)  affectedRows(%d) failed: %v", u, res.RowsAffected(), err)
	}
	return err
}

func (u *User) SetPassword(word string) error {
	res, err := Db.Model(u).
		Set("password = ?", word).
		Where("id = ?", u.ID).
		Update()
	if err != nil {
		log.Errorf("SetPassword for User(%v)  affectedRows(%d) failed: %v", u, res.RowsAffected(), err)
	}
	return err
}

func (u *User) SetAvatarURL(url string) error {
	res, err := Db.Model(u).
		Set("avatar_url = ?", url).
		Where("id = ?", u.ID).
		Update()
	if err != nil {
		log.Errorf("SetAvatarURL for User(%v)  affectedRows(%d) failed: %v", u, res.RowsAffected(), err)
	}
	return err
}

func CreateUser(email string, name string, password string) (*User, error) {
	var s = User{
		Name:       name,
		Email:      email,
		Password:   password,
		CreateTime: time.Now(),
		Role:       1,
		Disabled:   false,
		// 暂定 4G
		MaxSize: 1024 * 1024 * 1024 * 4,
		CurSize: 0,
	}
	err := Db.Insert(&s)
	return &s, err
}

func NewUser(form *User) error {
	return Db.Insert(form)
}
