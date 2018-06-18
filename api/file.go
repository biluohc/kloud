package api

import (
	"io"
	"net/url"
	"os"
	"strings"

	log "github.com/cihub/seelog"

	"github.com/biluohc/kloud/model"
	"github.com/labstack/echo"
)

type Sha1Form struct {
	Name string `json:"name" form:"name" query:"name"`
	Sha1 string `json:"sha1" form:"sha1" query:"sha1"`
}

func CreateFileBySha1(c echo.Context) (err error) {
	form := new(Sha1Form)
	err = nil
	err = c.Bind(form)

	if err != nil {
		return Err(4000, "表单错误")
	}

	pid := strings.TrimSpace(c.Param("id"))
	if pid == "" {
		return Err(4000, "父目录ID不能为空")
	}

	cc := c.(*KContext)
	user := cc.User()

	nameError := CheckEntryName(form.Name)
	if nameError != nil {
		return nameError
	}

	file, err := model.FindFileBySha1(strings.TrimSpace(form.Sha1))
	if err != nil {
		return Err(4040, "没有该文件")
	}

	if user.CurSize+file.Size >= user.MaxSize {
		return Err(4030, "当前用户容量已经用完")
	}

	entry, err := model.CreateEntryByUid(form.Name, pid, file.ID, user.ID)
	entry.File = file

	if err != nil {
		Err(500, err.Error())
	}
	_ = file.AddRc()
	_ = user.AddCurSize(file.Size)

	return Ok(200, entry)
}

func CreateFile(c echo.Context) (err error) {
	cc := c.(*KContext)
	user := cc.User()
	// Source
	pid := strings.TrimSpace(c.Param("id"))
	if pid == "" {
		return Err(4000, "父目录ID不能为空")
	}

	file, err := c.FormFile("file")

	if err != nil {
		return Err(4002, "上传文件表单错误")
	}
	src, err := file.Open()

	if file.Size < 0 {
		return Err(4002, "传文件表单错误: 上传文件大小<0")
	}
	if file.Size+user.CurSize >= user.MaxSize {
		return Err(4030, "当前用户容量已经用完")
	}
	if err == nil {
		defer src.Close()

		// 防止名字里有对服务器文件系统非法的字符
		tmp := ".tmp/" + model.NanoId() + "_" + url.PathEscape(file.Filename)
		dst, err := os.Create(tmp)

		log.Infof("%s: %v", tmp, err)
		if err == nil {
			bytes, err := io.Copy(dst, src)
			if err == nil || bytes != file.Size {
				f, err := model.CreateFile(file.Filename, tmp, bytes)
				if err == nil {
					// 以后要处理文件处理文件重名的问题，百度是在后面直接加(1)
					// 移动/保存分享时报错
					entry, err := model.CreateEntryByUid(file.Filename, pid, f.ID, user.ID)
					if err == nil {
						_ = user.AddCurSize(f.Size)
						entry.File = f
						return OkWithMsg(200, "上传文件成功", entry)
					} else {
						log.Criticalf("CreateEntryFailed AfterUploadFilesuccess: %s %v", file.Filename, err)
					}
				}
			} else {
				_ = dst.Close()
				e := os.Remove(tmp)
				// 曾看到上面创建的文件这里删不掉->只读，umask又是对的。。
				if e != nil {
					log.Errorf("Remove tmpFile(%s) failed: %s", tmp, e)
				}
			}
		}
	}
	return Err(400, "上传文件失败")
}
