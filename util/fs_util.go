package util

import (
	"os"
)

func Dir_not_exist(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return false
	} else {
		return true
	}
}

func File_not_exist(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() == false {
		return false
	} else {
		return true
	}
}
