package model

import (
	"os"
)

const LogPath string = "log/dist.txt"

type Stat struct {
	UsersNumber int   `json:"usersNumber"`
	FilesNumber int64 `json:"filesNumber"`
	LogSize     int64 `json:"logSize"`
}

func GetStat() (*Stat, error) {
	var stat Stat
	uc, err := Db.Model((*User)(nil)).Count()
	if err != nil {
		return &stat, err
	}
	fc, err := Db.Model((*File)(nil)).Count()
	if err != nil {
		return &stat, err
	}
	logInfo, err := os.Stat(LogPath)

	if err != nil {
		return &stat, err
	}

	stat.UsersNumber, stat.FilesNumber, stat.LogSize = uc, int64(fc), logInfo.Size()
	return &stat, err
}

func ClearLogFile() error {
	return os.Truncate(LogPath, 0)
}
