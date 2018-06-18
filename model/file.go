package model

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"time"

	log "github.com/cihub/seelog"
	"github.com/go-pg/pg/orm"
)

func sha1File(fName string) (string, error) {
	f, e := os.Open(fName)
	if e != nil {
		log.Error("sha1File.OpenFile: ", e)
		return "", e
	}
	h := sha1.New()
	_, e = io.Copy(h, f)
	if e != nil {
		log.Error("sha1FileIoCopy: ", e)
		return "", e
	}
	return hex.EncodeToString(h.Sum(nil)), e
}

type File struct {
	tableName  struct{}  `sql:"kloud_file"`
	ID         string    `sql:"id,pk" json:"id"`
	Sha1       string    `db:"sha1,notnull" json:"sha1"`
	Size       int64     `db:"size,notnull" json:"size"`
	Path       string    `db:"path,notnull" json:"path"`
	Rc         int32     `db:"rc,notnull" json:"rc"`
	Disabled   bool      `db:"disabled,notnull" json:"disabled"`
	Name       string    `db:"name,notnull" json:"name"`
	CreateTime time.Time `db:"create_time,notnull" json:"createTime" `
}

func (u *File) BeforeInsert(db orm.DB) error {
	if u.ID == "" {
		u.ID = NanoId()
	}
	return nil
}

func FindFileBySha1(sha1 string) (*File, error) {
	var u = new(File)
	err := Db.Model(u).
		Where("sha1 = ?", sha1).
		Select()
	return u, err
}

func FindFileByID(id string) (*File, error) {
	var u = new(File)
	err := Db.Model(u).
		Where("id = ?", id).
		Select()
	return u, err
}

func GetFiles() ([]File, error) {
	var files []File
	err := Db.Model(&files).
		Order("create_time desc").
		Select()
	return files, err
}

func (f *File) SetDisable(b bool) error {
	res, err := Db.Model(f).
		Set("disabled = ?", b).
		Where("id = ?", f.ID).
		Update()
	if err != nil {
		log.Errorf("SetDisable(%v) for File(%v)  affectedRows(%d) failed: %v", b, f, res.RowsAffected(), err)
	}
	return err
}

// update kloud_file as file set rc=file.rc+1 where id='3d60b3b0-bc32-4581-a68c-3d92e173b8dc';
func (f *File) SubRcAndRmIfZero() error {
	var rc int

	res, err := Db.Model(f).
		Set("rc = rc-1").
		Where("id = ?", f.ID).
		Returning("rc").
		Update(&rc)

	f.Rc = int32(rc)
	if err != nil {
		log.Errorf("SubRc for File(%v)  affectedRows(%d) failed: %v", f, res.RowsAffected(), err)
		return err
	}

	log.Info(rc, f)
	// 这些逻辑并不会影响当前api
	if f.Rc == 0 {
		// 防止在删除时又有请求把Rc加上了。。
		res, err := Db.Model((*File)(nil)).Where("id = ? and rc = 0", f.ID).Delete()
		rows := res.RowsAffected()
		if rows > 1 {
			log.Criticalf("SubRcAndRmIfZero_RmRows %d>1: %v", rows, err)
		}
		if rows > 0 {
			e := os.Remove(f.Path)
			if e != nil {
				log.Criticalf("Remove %s failed: %s", f.Path, e)
			}
		}
	}
	return nil
}

func (f *File) AddRc() error {
	res, err := Db.Model(f).
		Set("rc = rc+1").
		Where("id = ?", f.ID).
		Update()
	if err != nil {
		log.Errorf("AddRc for File(%v)  affectedRows(%d) failed: %v", f, res.RowsAffected(), err)
	}
	return err
}

func CreateFile(name string, path string, size int64) (*File, error) {
	sha1, err := sha1File(path)
	if err != nil {
		log.Criticalf("sha1File(%s) failed: %v", path, err)
	}

	fy, err := FindFileBySha1(sha1)
	// 已经有了，不考虑 sha1 碰撞
	if err == nil {
		e := os.Remove(path)
		log.Criticalf("Remove TmpFile failed: %s %v", path, e)
		_ = fy.AddRc()
		return fy, err
	}

	// 我只能说， 蠢的一逼。
	// 月份 1,01,Jan,January
	// 日　 2,02,_2
	// 时　 3,03,15,PM,pm,AM,am
	// 分　 4,04
	// 秒　 5,05
	// 年　 06,2006
	newDir := time.Now().Local().Format(".fs/2006-01/02/")

	// 这如果有了并不会失败
	err = os.MkdirAll(newDir, 0755)
	if err != nil {
		log.Criticalf("os.MkdirAll(%s, 0755) failed: %v", newDir, err)
		return nil, err
	}

	newPath := newDir + sha1
	err = os.Rename(path, newPath)

	if err != nil {
		log.Criticalf("Move tmpFile to Fs failed: %s -> %s : %s", path, newPath, err)
		return nil, err
	}
	var f = File{
		Path:     newPath,
		Sha1:     sha1,
		Size:     size,
		Name:     name,
		Disabled: false,
		Rc:       1,
	}

	err = Db.Insert(&f)
	return &f, err
}
