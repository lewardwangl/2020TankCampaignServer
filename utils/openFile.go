package utils

import (
	"os"
	"path/filepath"
)

func OpenFile(path string) (f *os.File, err error) {
	f, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		dir, _ := filepath.Split(path)
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			panic("这就出了大问题了，" + err.Error())
		}
		f, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}
	return
}
