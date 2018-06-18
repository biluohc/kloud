package api

import (
	"strings"
	"time"

	"github.com/biluohc/kloud/model"
	log "github.com/cihub/seelog"
	"github.com/labstack/echo"
)

func GetUsers(c echo.Context) (err error) {
	users, err := model.GetUsers()
	if err != nil {
		return Err(500, err.Error())
	}
	return Ok(200, users)
}

func GetUser(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}

	user, err := model.FindUserByID(id)
	if err != nil {
		return Err(500, err.Error())
	}
	return Ok(200, user)
}

type UserForm struct {
	Name      string `json:"name" form:"name" query:"name"`
	Email     string `json:"email" form:"email" query:"email"`
	Password  string `json:"password" form:"password" query:"password"`
	AvatarURL string `json:"avatarURL"`
	MaxSize   int64  `json:"maxSize"`
	Role      int    `json:"role"`
}

func UpdateAvatarURL(c echo.Context) (err error) {
	cc := c.(*KContext)
	user := cc.User()

	var form = new(UserForm)
	err = c.Bind(form)
	if err != nil || form.AvatarURL == "" {
		return Err(4000, "表单错误")
	}
	_ = user.SetAvatarURL(form.AvatarURL)
	return Ok(200, user)
}

func UpdatePassword(c echo.Context) (err error) {
	cc := c.(*KContext)
	user := cc.User()

	var form = new(UserForm)
	err = c.Bind(form)
	if err != nil || form.Password == "" {
		return Err(4000, "表单错误")
	}
	_ = user.SetPassword(form.Password)
	return Ok(200, user)
}

func NewUser(c echo.Context) (err error) {
	var form = new(UserForm)
	err = c.Bind(form)
	if err != nil {
		return Err(4000, "表单错误")
	}

	form.Name = strings.TrimSpace(form.Name)
	form.Email = strings.TrimSpace(form.Email)
	form.Password = strings.TrimSpace(form.Password)

	if form.Name == "" || form.Email == "" || form.Password == "" {
		return Err(4000, "用户名、邮箱和密码不能为空")
	}

	if form.MaxSize < 0 {
		return Err(4000, "总容量不能小于0")
	}
	if form.Role < 0 || form.Role > 1 {
		return Err(4000, "非法角色 0<||>1")
	}

	var user = model.User{
		Name:       form.Name,
		Email:      form.Email,
		Password:   form.Password,
		CreateTime: time.Now(),
		Role:       uint16(form.Role),
		MaxSize:    form.MaxSize,
		CurSize:    0,
		Disabled:   false,
	}
	err = model.NewUser(&user)
	if err != nil {
		return Err(400, err.Error())
	}
	return Ok(200, user)
}

func ModifyUser(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}

	var form = new(UserForm)
	err = c.Bind(form)
	if err != nil {
		return Err(4000, "表单错误")
	}

	form.Name = strings.TrimSpace(form.Name)
	form.Email = strings.TrimSpace(form.Email)
	form.Password = strings.TrimSpace(form.Password)

	if form.Name == "" || form.Email == "" {
		return Err(4000, "用户名和邮箱不能为空")
	}

	if form.MaxSize < 0 {
		return Err(4000, "总容量不能小于0")
	}
	if form.Role < 0 || form.Role > 1 {
		return Err(4000, "非法角色 0<||>1")
	}

	user, err := model.FindUserByID(id)
	if err != nil {
		return Err(4040, "没有该用户")
	}

	user.Name = form.Name
	user.Email = form.Email
	user.Role = uint16(form.Role)
	user.AvatarURL = form.AvatarURL
	user.MaxSize = form.MaxSize
	// 密码有才更新
	if form.Password != "" {
		user.Password = form.Password
	}

	err = user.Update()
	if err != nil {
		return Err(500, err.Error())
	}
	return Ok(200, user)
}

func UserDisable(c echo.Context) (err error) {
	log.Info("disable")
	return setUserDisable(c, true)
}

func UserUnDisable(c echo.Context) (err error) {
	log.Info("undisable")
	return setUserDisable(c, false)
}

func setUserDisable(c echo.Context, b bool) error {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}
	user, err := model.FindUserByID(id)
	if err != nil {
		return Err(4040, "没有该用户")
	}
	err = user.SetDisabled(b)
	if err != nil {
		return Err(500, err.Error())
	}
	return Ok(200, user)
}

// SignUp
type SignUpForm struct {
	Name      string `db:"name,notnull" json:"name"`
	Email     string `json:"email" form:"email" query:"email"`
	Password  string `json:"password" form:"password" query:"password"`
	AvatarURL string `db:"avatar_url" json:"avatarURL"`
}

func Signin(c echo.Context) (err error) {
	u := new(SignUpForm)
	if err = c.Bind(u); err != nil {
		return Err(4000, "表单错误")
	}

	user, err := model.CreateUser(u.Email, u.Name, u.Password)
	if err != nil {
		return Err(400, "用户Emai或名字已经被占用")
	}
	return Ok(200, user)
}
