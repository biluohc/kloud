package api

import (
	"strings"

	"github.com/biluohc/kloud/model"
	"github.com/labstack/echo"
)

func CheckEntryName(name string) *Result {
	if name == "" {
		return Err(4000, "文件/目录名字不能为空")
	}

	if name == "." || name == ".." {
		return Err(4000, "文件/目录名字不能为`.`或'..'")
	}

	if strings.Contains(name, "/") {
		return Err(4000, "文件/目录名字不能含`/`")
	}
	return nil
}

type NameForm struct {
	Name string `json:"name" form:"name" query:"name"`
}

func CreateDir(c echo.Context) (err error) {
	form := new(NameForm)
	err = nil
	err = c.Bind(form)
	form.Name = strings.TrimSpace(form.Name)

	if err != nil {
		return Err(4000, "表单错误")
	}

	pid := strings.TrimSpace(c.Param("id"))
	if pid == "" {
		return Err(4000, "父目录ID不能为空")
	}

	nameError := CheckEntryName(form.Name)
	if nameError != nil {
		return nameError
	}

	cc := c.(*KContext)
	entry, err := model.CreateEntryByUid(form.Name, pid, "", cc.User().ID)
	if err != nil {
		return Err(500, err.Error())
	}
	return OkWithMsg(200, "创建目录成功", entry)
}
