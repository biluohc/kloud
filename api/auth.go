package api

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"

	"github.com/biluohc/kloud/model"
)

func AdminAuthRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*KContext)
		if cc.User().Role != 0 {
			return Err(4032, "您的帐号不是管理员")
		}
		return next(c)
	}
}
func AuthRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)

		s401 := err != nil || sess.Values["id"] == nil

		if !s401 {
			if id, ok := sess.Values["id"].(string); ok {
				se, err := model.FindSessByID(id)
				if err != nil {
					s401 = true
				} else {
					user, err := se.User()
					if err != nil {
						s401 = true
					} else if user.Disabled {
						// user being disabled
						return Err(4013, "您的帐号被禁用，请联系管理员")
					} else {
						cc := c.(*KContext)
						cc.SetUser(user)
					}
				}
			}
		}

		if s401 {
			return Err(401, "请用合法帐号登录")
		}

		return next(c)
	}
}

func GetUndisabledUserFromSession(c echo.Context) *model.User {
	sess, err := session.Get("session", c)
	if err != nil || sess.Values["id"] == nil {
		return nil
	}
	if id, ok := sess.Values["id"].(string); ok {
		se, err := model.FindSessByID(id)
		if err == nil {
			user, err := se.User()
			if err == nil && !user.Disabled {
				return user
			}
		}
	}
	return nil
}

// 返回当前登录用户信息
func CurLogin(c echo.Context) (err error) {
	cc := c.(*KContext)
	return Ok(200, cc.User())
}

// 退出登录
func ExitLogin(c echo.Context) (err error) {
	sess, _ := session.Get("session", c)
	id, _ := sess.Values["id"].(string)
	err = model.DeleteSessByID(id)
	if err != nil {
		return OkWithMsg(200, "退出登录成功", "")
	}
	return Err(404, "找不到该Session")
}

// Login
type LoginForm struct {
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
	// 最大68年，完全没问题
	MaxAge int `json:"maxAge" form:"maxAge" query:"maxAge"`
}

func Login(c echo.Context) (err error) {
	u := new(LoginForm)
	if err = c.Bind(u); err != nil {
		return Err(4000, "表单错误")
	}

	user, err := model.FindUserByEmail(u.Email)
	if err != nil {
		return Err(4011, "帐号不存在")
	}

	if user.Password == u.Password {
		sess, err := session.Get("session", c)
		hassess := err == nil && sess.Values["id"] != nil
		elsetry := false

		if hassess {
			if id, ok := sess.Values["id"].(string); ok {
				se, err := model.FindSessByID(id)
				if err != nil || se.UserID != user.ID {
					elsetry = true
				} else {
					return OkWithMsg(200, "已经登录", user)
				}
			} else {
				elsetry = true
			}
		}

		if !hassess || elsetry {
			ua := c.Request().Header.Get("User-Agent")
			ip := c.RealIP()
			se, err := model.CreateSess(user.ID, u.MaxAge, ua, ip)
			if err != nil {
				return Err(500, "创建Session失败")
			}
			// 86400 = 1 day
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   u.MaxAge,
				HttpOnly: true,
			}
			sess.Values["id"] = se.ID
			sess.Save(c.Request(), c.Response())
			return OkWithMsg(200, "登录成功", user)
		}
		// 貌似不会走到此分支
		return Err(401, "未知错误")
	} else {
		return Err(4012, "帐号或者密码错误")
	}
}
