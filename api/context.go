package api

import (
	"github.com/labstack/echo"

	"github.com/biluohc/kloud/model"
)

type KContext struct {
	echo.Context
	user *model.User
}

func (c *KContext) SetUser(u *model.User) {
	c.user = u
}

func (c *KContext) User() *model.User {
	return c.user
}

func KloudContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &KContext{c, nil}
		return next(cc)
	}
}
