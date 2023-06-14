package controller

import "os"

func rcloneConfigFile() string {
	filePath := "/root/.config/rclone/rclone.conf"
	if _, err := os.Stat(filePath); err != nil {
		return ""
	}
	return filePath
}
