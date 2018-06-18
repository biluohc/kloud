package model

import (
	log "github.com/cihub/seelog"
)

func (e *Entry) DeleteByUser(user *User) error {
	// 目录
	if e.FileID == "" {
		entris, err := FindEntrySubsByIDandUid(e.ID, user.ID, false)
		if err != nil {
			return err
		}
		for i := range entris {
			err = entris[i].DeleteByUser(user)
			if err != nil {
				return err
			}
		}
		// 删除目录本身
		res, err := Db.Model(e).Where("id = ? and user_id = ?", e.ID, user.ID).Delete()
		if res.RowsAffected() != 1 {
			log.Criticalf("Delete Dir(%v) failed affactedRows %d: %v", e, res.RowsAffected(), err)
		}
		if err != nil {
			return err
		}
	} else {
		// 文件
		f, err := FindFileByID(e.FileID)
		if err != nil {
			return err
		}

		res, err := Db.Model(e).Where("id = ? and user_id = ?", e.ID, user.ID).Delete()
		if res.RowsAffected() != 1 {
			log.Criticalf("Delete File(%v) failed affactedRows %d: %v", e, res.RowsAffected(), err)
		}

		err = f.SubRcAndRmIfZero()
		if err == nil {
			_ = user.SubCurSize(f.Size)
		}
		return err
	}
	return nil
}
