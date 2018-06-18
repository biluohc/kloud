// 拦截一切referer, 防止盗图之类

package api

import (
	"crypto/sha256"
	"fmt"
	"mime"
	"path"
	"strings"

	"github.com/biluohc/kloud/model"
	log "github.com/cihub/seelog"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

type ShareForm struct {
	Password string `json:"password" form:"password" query:"password"`
	// 有效时间 s ，由前端组合 1/7日1月1年
	MaxAge int `json:"maxAge" form:"maxAge" query:"maxAge"`
}

func EntryShare(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}

	var form = new(ShareForm)
	// 没有的话，会爆炸吗？
	err = c.Bind(form)
	if err != nil {
		return Err(4000, "表单错误")
	}

	cc := c.(*KContext)
	entry, err := model.FindEntryByIDandUid(id, cc.User().ID, false)
	if err != nil {
		return Err(4041, "找不到该条目")
	}

	s, _ := entry.ShareWithPassWordAndMaxAge(form.Password, form.MaxAge)

	return Ok(200, s)
}

// 修改密码
func SharePassword(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}
	var form = new(ShareForm)
	// 没有的话，会爆炸吗？
	err = c.Bind(form)
	if err != nil {
		return Err(4000, "表单错误")
	}

	cc := c.(*KContext)
	share, err := model.FindShareByIDandUid(id, cc.User().ID)
	if err != nil {
		return Err(4041, "找不到该分享")
	}
	share.SetPassword(form.Password)
	return Ok(200, share)
}

func ShareMaxAge(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}
	var form = new(ShareForm)
	// 没有的话，会爆炸吗？
	err = c.Bind(form)
	if err != nil {
		return Err(4000, "表单错误")
	}

	cc := c.(*KContext)
	share, err := model.FindShareByIDandUid(id, cc.User().ID)
	if err != nil {
		return Err(4041, "找不到该分享")
	}
	share.SetMaxAge(form.MaxAge)
	return Ok(200, share)
}

func ShareDelete(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}
	cc := c.(*KContext)
	share, err := model.FindShareByIDandUid(id, cc.User().ID)
	if err != nil {
		return Err(4041, "找不到该条目")
	}
	err = share.Delete()
	if err != nil {
		return Err(500, err.Error())
	}
	return Ok(200, share)
}

func Shares(c echo.Context) (err error) {
	cc := c.(*KContext)
	shares, err := cc.User().Shares()
	if err != nil {
		return Err(4040, err.Error())
	}

	return Ok(200, shares)
}

func ShareInit(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}

	var form = new(ShareForm)
	// 没有的话，会爆炸吗？
	err = c.Bind(form)
	if err != nil {
		return Err(4000, "表单错误")
	}

	share, err := model.FindShareByID(id)

	if err != nil {
		return Err(4041, "找不到该条目")
	}

	if share.Password == "" {
		return OkWithMsg(200, "该分享没有密码", share)
	}

	curUser := GetUndisabledUserFromSession(c)
	if curUser != nil && curUser.ID == share.UserID {
		return OkWithMsg(200, "分享者即当前用户", share)
	}

	if form.Password == share.Password {
		passwordSha := passwordSha256(share.ID, share.Password)

		sess, err := session.Get("session", c)
		if err != nil {
			log.Error("浏览器不支持或关闭了Cookie", err)
			return Err(403, "获取Cookie失败，打开Cookie或使用支持Cookie的现代浏览器")
		}
		// age= 浏览器进程
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		}
		sess.Values[shareCookieId(share.ID)] = passwordSha
		sess.Save(c.Request(), c.Response())
		// js 得到 200 设置跳转到内容
		return OkWithMsg(200, "密码正确", share)
	}

	return Err(403, "密码错误")
}

func GetShareInfo(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4001, "非法Id")
	}

	share, err := model.FindShareByID(id)
	if err != nil {
		return Err(4041, "找不到该条目")
	}

	entry, _ := model.FindEntryByID(share.EntryID, true)
	share.Entry = entry

	share.Password = ""
	return Ok(200, share)
}

// 分享者是当前登录用户就可以直接访问
func GetShareFile(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4001, "非法Id")
	}

	referer := c.Request().Referer()
	if referer != "" && strings.Index(referer, Host) == -1 {
		return Err(4033, "Shared deny")
	}

	share, err := model.FindShareByID(id)
	if err != nil {
		return Err(4041, "找不到该条目")
	}

	curUser := GetUndisabledUserFromSession(c)
	// 分享者即当前用户
	curUserIsShareOwner := curUser != nil && curUser.ID == share.UserID

	if !curUserIsShareOwner && !hasRightPassword(c, share.ID, share.Password) {
		// 重定向到前端
		initURL := fmt.Sprintf("/#/share/%s/init", share.ID)
		return c.Redirect(302, initURL)
	}

	entry, err := model.FindEntryByID(share.EntryID, true)
	if err != nil {
		return Err(4041, "找不到该条目")
	}
	if entry.FileID == "" {
		return Err(4003, "该条目不是文件")
	}

	if entry.File.Disabled {
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

func hasRightPassword(c echo.Context, id, word string) bool {
	if word == "" {
		return true
	}
	sess, err := session.Get("session", c)
	if err == nil {
		if tmp, ok := sess.Values[shareCookieId(id)].(string); ok {
			wordSha := passwordSha256(id, word)
			if tmp == wordSha {
				return true
			}
		}
	}
	return false
}

func passwordSha256(id, word string) string {
	h := sha256.New()
	h.Write([]byte(id + word))
	return string(h.Sum(nil))
}

func shareCookieId(id string) string {
	return "share:" + id
}
