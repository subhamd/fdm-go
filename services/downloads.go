package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/subhamd/fdm-go/entities"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func DownloadFiles(downloadType entities.DownloadType, files []string) (entities.DownloadStatus, string) {
	downloadEntity := new(entities.DownloadEntity)
	id := uuid.New().String()
	downloadEntity.Id = id
	downloadEntity.DownloadType = downloadType
	downloadEntity.Files = files

	entities.DownloadIdToDownloadStatusMap[id] = downloadEntity

	var downloadStatus entities.DownloadStatus

	downloadEntity.Status = entities.DownloadStatusQueued
	switch downloadType {
	case entities.SerialDownload:
		downloadEntity.StartTime = time.Now()
		downloadStatus = downloadSerially(downloadEntity)
		downloadEntity.EndTime = time.Now()
	case entities.ConcurrentDownload:
		downloadStatus = downloadConcurrently(downloadEntity)
	}

	downloadEntity.Status = downloadStatus

	return downloadStatus, id
}

func downloadSerially(downloadEntity *entities.DownloadEntity) entities.DownloadStatus {
	for _, url := range downloadEntity.Files {
		err := downloadAndSaveFileLocally(url)
		if err != nil {
			log.Println(err)
			return entities.DownloadStatusFailed
		}
	}

	return entities.DownloadStatusSuccess
}

func downloadConcurrently(downloadEntity *entities.DownloadEntity) entities.DownloadStatus {

	return entities.DownloadStatusQueued
}



func downloadAndSaveFileLocally(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileNameWithExtension := path.Base(url)
	extension := filepath.Ext(fileNameWithExtension)
	fileName := strings.TrimSuffix(fileNameWithExtension, extension)

	fileNameWithExtension = fileName + uuid.New().String() + extension
	filePath := entities.DownloadFilePath + fileNameWithExtension

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	entities.DownloadedFilesToPathMap[fileNameWithExtension] = filePath
	return nil
}

func GetDownloadStatus(downloadId string) (*entities.DownloadEntity, error) {
	downloadEntity, keyExists := entities.DownloadIdToDownloadStatusMap[downloadId]

	if !keyExists {
		return nil, errors.New("unknown download ID")
	}

	return downloadEntity, nil
}
