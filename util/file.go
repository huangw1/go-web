package util

import "os"

func IsDir(filename string) bool {
	return isFileOrDir(filename, true)
}

func IsFile(filename string) bool {
	return isFileOrDir(filename, false)
}

func isFileOrDir(filename string, decideDir bool) bool {
	stat, err := os.Stat(filename)
	if err != nil {
		return false
	}
	isDir := stat.IsDir()
	if decideDir {
		return isDir
	}
	return !isDir
}