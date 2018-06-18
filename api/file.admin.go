package api

import (
	"mime"
	"path"
	"strings"

	"github.com/biluohc/kloud/model"

	"github.com/labstack/echo"
)

func GetFiles(c echo.Context) (err error) {
	files, err := model.GetFiles()
	if err != nil {
		return Err(500, err.Error())
	}
	return Ok(200, files)
}

func GetFile(c echo.Context) (err error) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4001, "非法Id")
	}

	file, err := model.FindFileByID(id)

	if err != nil {
		return Err(4040, "未找到该文件")
	}

	mime := mime.TypeByExtension(path.Ext(file.Name))

	if mime == "" {
		mime = "application/octet-stream"
	}
	header := c.Response().Header()
	header.Set("Content-Type", mime)
	return c.File(file.Path)
}

func GetStat(c echo.Context) (err error) {
	stat, err := model.GetStat()
	if err != nil {
		return Err(500, err.Error())
	}
	return Ok(200, stat)
}

func ClearLogFile(c echo.Context) error {
	return model.ClearLogFile()
}

func FileDisable(c echo.Context) (err error) {
	return setFileDisable(c, true)
}

func FileUnDisable(c echo.Context) (err error) {
	return setFileDisable(c, false)
}

func setFileDisable(c echo.Context, b bool) error {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return Err(4000, "ID不能为空")
	}
	f, err := model.FindFileByID(id)
	if err != nil {
		return Err(404, "找不到该文件")
	}
	err = f.SetDisable(b)
	if err != nil {
		return Err(500, err.Error())
	}
	return Ok(200, f)
}
