package api

import (
	"mime"
	"path"
	"strings"

	log "github.com/cihub/seelog"

	"github.com/biluohc/kloud/model"
	"github.com/labstack/echo"
)

func GetEntryInfo(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4001, "非法Id")
	}

	cc := c.(*KContext)
	entry, err := model.FindEntryByIDandUid(id, cc.User().ID, true)
	if err != nil {
		return Err(4041, "找不到该条目")
	}
	return Ok(200, entry)
}

// For the Entry is file
func GetEntryFile(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4001, "非法Id")
	}

	cc := c.(*KContext)
	entry, err := model.FindEntryByIDandUid(id, cc.User().ID, true)
	if err != nil {
		log.Info(entry, err)
		return Err(4041, "找不到该条目")
	}

	if entry.FileID == "" {
		return Err(4003, "该条目不是文件")
	}

	if entry.File.Disabled && cc.User().Role > 0 {
		return Err(4031, "该文件被禁用，请联系管理员")
	}
	mime := mime.TypeByExtension(path.Ext(entry.Name))

	if mime == "" {
		mime = "application/octet-stream"
	}
	header := c.Response().Header()
	header.Set("Content-Type", mime)
	return c.File(entry.File.Path)
}

func GetEntryRootInfo(c echo.Context) (err error) {
	cc := c.(*KContext)
	entry, err := model.FindOrCreateRootByUid(cc.User())
	if entry.FileID != "" {
		log.Criticalf("%s's Root(%s) is File(%s)", cc.User(), entry.ID, entry.FileID)
	}
	return Ok(200, entry)
}

func GetEntryRootPath(c echo.Context) (err error) {
	cc := c.(*KContext)
	entry, err := model.FindOrCreateRootByUid(cc.User())
	if entry.FileID != "" {
		log.Criticalf("%s's Root(%s) is File(%s)", cc.User(), entry.ID, entry.FileID)
	}
	entries := make([]model.Entry, 1)
	entries[0] = *entry

	return Ok(200, entries)
}

func GetEntryPath(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	cc := c.(*KContext)
	entry, err := model.FindEntryByIDandUid(id, cc.User().ID, false)
	if err != nil {
		return Err(4041, "找不到该条目")
	}

	es, err := model.FindEntryPath(entry.ID)
	if err != nil {
		return Err(500, err.Error())
	}
	tmp := make([]model.Entry, len(es))
	// 反转
	for i := 0; i < len(es); i++ {
		tmp[i] = es[len(es)-1-i]
	}

	return Ok(200, tmp)
}

func getEntrySubs(c echo.Context, id string) (err error) {
	cc := c.(*KContext)

	entries, err := model.FindEntrySubsByIDandUid(id, cc.User().ID, true)

	log.Info(id, cc.User().ID, len(entries), err)
	if err != nil {
		return Err(404, err.Error())
	}

	// 坑
	if entries == nil {
		entries = make([]model.Entry, 0)
	}

	return Ok(200, entries)
}

func GetEntryRootSubs(c echo.Context) (err error) {
	cc := c.(*KContext)
	id, err := model.FindOrCreateRootByUid(cc.User())
	if err != nil {
		return Err(500, err.Error())
	}
	return getEntrySubs(c, id.ID)
}

func GetEntrySubs(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4001, "非法ID")
	}
	return getEntrySubs(c, id)
}

func GetEntryTrashes(c echo.Context) (err error) {
	cc := c.(*KContext)
	entries, err := model.FindTrashesByUid(cc.User().ID)
	if err != nil {
		return Err(404, err.Error())
	}
	return Ok(200, entries)
}
