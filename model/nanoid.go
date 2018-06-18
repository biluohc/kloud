package model

import (
	"github.com/matoous/go-nanoid"
)

var table string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+-"

// https://github.com/nikolay-govorov/nanoid
// 这个Go版的Nanoid()生成的位数默认是22位, 稍微换了编码集
// A-Za-z0-9_~ -> A-Za-z0-9+_
func NanoId() string {
	id, _ := gonanoid.Generate(table, 22)
	return id
}
