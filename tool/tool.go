package tool

import (
	"errors"
	"os"
	"regexp"
)

// PathExists check if path exists
func PathExists(path string) error {
	_, err := os.Stat(path)
	return err
}

// MoveFile move file from to
func MoveFile(from, to string) error {
	return os.Rename(from, to)
}

// GetDumpFileName get dump file name from stdout
func GetDumpFileName(str string) string {
	re := regexp.MustCompile(`gitea-dump-\d+.zip`)
	return re.FindString(str)
}

// DeleteLocal delete the oldest file if there are more than maxNum
func DeleteLocal(path string, maxNum int) error {
	// 读取 path 下的文件列表
	if fileList, err := os.ReadDir(path); errors.Is(err, nil) {
		for _, file := range fileList {
			if getDirFileNum(path) <= maxNum {
				break
			}
			_ = os.Remove(path + file.Name())
		}
	} else {
		return err
	}
	return nil
}

func getDirFileNum(path string) int {
	if list, err := os.ReadDir(path); errors.Is(err, nil) {
		return len(list)
	}
	return 0
}
