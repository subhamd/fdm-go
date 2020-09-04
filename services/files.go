package services

import (
	"github.com/subhamd/fdm-go/entities"
)

func getFiles() {
	for key, value := range entities.DownloadedFilesToPathMap {
		print(key, value)
	}
}
