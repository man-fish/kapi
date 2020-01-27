package utils

import (
	"os"
	"path/filepath"
)

func RootPath() string {
	workpath, _ := os.Executable()
	rootPath := filepath.Join(filepath.Dir(workpath),"./../../")
	return rootPath
}
/* 获取项目根路径 */
