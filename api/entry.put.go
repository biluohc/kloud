package api

import (
	"fmt"
	"strings"

	"github.com/biluohc/kloud/model"
	log "github.com/cihub/seelog"
	"github.com/labstack/echo"
)

func EntryRename(c echo.Context) (err error) {
	form := new(NameForm)
	err = nil
	err = c.Bind(form)

	if err != nil {
		return Err(4000, "表单错误")
	}

	name := strings.TrimSpace(form.Name)
	nameError := CheckEntryName(form.Name)
	if nameError != nil {
		return nameError
	}

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}

	cc := c.(*KContext)
	it, res, err := model.RenameEntryByIDandUid(name, id, cc.User().ID)
	if res.RowsAffected() == 1 {
		return OkWithMsg(200, "重命名成功", it)
	} else if res.RowsAffected() > 1 {
		log.Errorf("EntryRename RowsAffected %d > 1: %v  %v", res.RowsAffected(), form, err)
	}

	return Err(400, "重命名失败")
}

type MoveForm struct {
	ID string `json:"id" form:"id" query:"id"`
}

func EntryMoveTo(c echo.Context) (err error) {
	form := new(MoveForm)
	err = nil
	err = c.Bind(form)

	if err != nil {
		return Err(4000, "表单错误")
	}

	if form.ID == "" {
		return Err(4000, "源ID不能为空")
	}

	pid := strings.TrimSpace(c.Param("id"))
	if pid == "" {
		return Err(4000, "目标ID不能为空")
	}

	// 这个还不能过滤移动到自己子目录里，全部由前端负责吧，223
	if pid == form.ID {
		return Err(4000, "目标ID不能为源ID")
	}

	cc := c.(*KContext)

	entry, res, err := model.MoveEntryByIDandUid("", pid, form.ID, cc.User().ID)
	if res.RowsAffected() != 1 {
		log.Errorf("EntryMove RowsAffected %d !=1: %v %v", res.RowsAffected(), form, err)
	}

	if err == nil {
		return OkWithMsg(200, fmt.Sprintf("%s移动成功", entry.Name), entry)
	}
	return Err(400, fmt.Sprintf("%s移动失败", entry.Name))
}

func EntryTrash(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}

	cc := c.(*KContext)
	it, res, err := model.TrashbyIDandUid(id, cc.User().ID)
	if res.RowsAffected() == 1 {
		return OkWithMsg(200, "移动到回收站成功", it)
	} else if res.RowsAffected() > 1 {
		log.Errorf("EntryTrash RowsAffected %d > 1: %v %v", res.RowsAffected(), id, err)
	}

	return Err(400, "移动到回收站失败")
}

func EntryUnTrash(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}
	cc := c.(*KContext)
	it, res, err := model.UnTrashbyIDandUid(id, cc.User().ID)

	log.Info(it, res, err)
	if res.RowsAffected() == 1 {
		return OkWithMsg(200, "从回收站恢复成功", it)
	} else if res.RowsAffected() > 1 {
		log.Errorf("EntryUnTrash RowsAffected %d > 1: %v %v", res.RowsAffected(), it, err)
	}

	return Err(400, "从回收站恢复失败")
}
