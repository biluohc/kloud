package api

import (
	"strings"

	"github.com/biluohc/kloud/model"
	"github.com/labstack/echo"
)

func EntryDelete(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}

	cc := c.(*KContext)
	entry, err := model.FindEntryByIDandUid(id, cc.User().ID, false)

	if err != nil {
		return Err(4041, "找不到该条目")
	}

	err = entry.DeleteByUser(cc.User())
	if err != nil {
		return Err(500, "删除条目错误")
	}
	return Ok(200, "")
}
