package api

// 拦截域名referer, 防止盗图之类， 如果是ip可能是在本地调试。

import (
	"mime"
	"path"
	"regexp"
	"strings"

	"github.com/biluohc/kloud/model"
	log "github.com/cihub/seelog"
	"github.com/labstack/echo"
)

type PublicForm struct {
	Referer string `json:"referer" form:"referer" query:"referer"`
	// 有效时间 s ，由前端组合
	MaxAge int `json:"maxAge" form:"maxAge" query:"maxAge"`
}

func EntryPublic(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}

	var form = new(PublicForm)
	// 没有的话，会爆炸吗？
	err = c.Bind(form)
	if err != nil {
		return Err(4000, "表单错误")
	}

	_, err = regexp.Compile(form.Referer)
	if err != nil {
		return Err(4004, "正则错误")
	}

	// 空的字符串和 .* 一样的效果
	cc := c.(*KContext)
	entry, err := model.FindEntryByIDandUid(id, cc.User().ID, false)
	if err != nil {
		return Err(4041, "找不到该条目")
	}

	s, _ := entry.PublicWithReferers(form.Referer, form.MaxAge)

	return Ok(200, s)
}

func PublicReferer(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}
	var form = new(PublicForm)
	// 没有的话，会爆炸吗？
	err = c.Bind(form)
	if err != nil {
		return Err(4000, "表单错误")
	}

	// 检测合法性
	_, err = regexp.Compile(form.Referer)
	if err != nil {
		return Err(4004, "正则错误")
	}

	cc := c.(*KContext)
	public, err := model.FindPublicByIDandUid(id, cc.User().ID)
	if err != nil {
		return Err(4041, "找不到该分享")
	}

	_ = public.SetReferer(form.Referer)

	return Ok(200, public)
}

func PublicMaxAge(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}
	var form = new(PublicForm)
	// 没有的话，会爆炸吗？
	err = c.Bind(form)
	if err != nil {
		return Err(4000, "表单错误")
	}

	cc := c.(*KContext)
	share, err := model.FindPublicByIDandUid(id, cc.User().ID)
	if err != nil {
		return Err(4041, "找不到该分享")
	}
	share.SetMaxAge(form.MaxAge)
	return Ok(200, share)
}

func PublicDelete(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}
	cc := c.(*KContext)
	share, err := model.FindPublicByIDandUid(id, cc.User().ID)
	if err != nil {
		return Err(4041, "找不到该条目")
	}
	err = share.Delete()
	if err != nil {
		return Err(500, err.Error())
	}
	return Ok(200, share)
}

func Publics(c echo.Context) (err error) {
	cc := c.(*KContext)
	shares, err := cc.User().Publics()
	if err != nil {
		return Err(4040, err.Error())
	}
	return Ok(200, shares)
}

// 检测referer
func GetPublicFile(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4001, "非法Id")
	}

	log.Info(id)

	referer := c.Request().Referer()

	public, err := model.FindPublicByID(id)
	if err != nil {
		return Err(4041, "找不到该条目")
	}

	if !RefererValid(public.Referer, referer) {
		return Err(403, "")
	}

	entry, err := model.FindEntryByID(public.EntryID, false)
	if err != nil {
		return Err(4041, "找不到该条目")
	}
	if entry.FileID == "" {
		return Err(4003, "该条目不是文件")
	}

	file, err := model.FindFileByID(entry.FileID)

	if file.Disabled {
		return Err(4031, "该文件被禁用，请联系管理员")
	}
	mime := mime.TypeByExtension(path.Ext(entry.Name))

	if mime == "" {
		mime = "application/octet-stream"
	}
	header := c.Response().Header()
	header.Set("Content-Type", mime)
	return c.File(file.Path)
}

func RefererValid(allows, referer string) bool {
	if allows == "" || referer == "" || strings.HasPrefix(referer, Host) {
		return true
	}
	b, err := regexp.MatchString(allows, referer)
	if err != nil {
		log.Warnf("RefererValid正则错误: %v", err)
	}
	return b
}
