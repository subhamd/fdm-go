package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/subhamd/fdm-go/entities"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
		downloadStatus = downloadSerially(downloadEntity) // sync
		downloadEntity.EndTime = time.Now()
	case entities.ConcurrentDownload:
		go downloadConcurrently(downloadEntity) // async
		downloadStatus = entities.DownloadStatusQueued
	}

	downloadEntity.Status = downloadStatus

	return downloadStatus, id
}

func downloadSerially(downloadEntity *entities.DownloadEntity) entities.DownloadStatus {
	baseFilePath, err := createAndGetBaseFilePath(downloadEntity)
	if err != nil {
		log.Println(err)
		return entities.DownloadStatusFailed
	}

	for _, url := range downloadEntity.Files {
		err := downloadAndSaveFileLocally(url, baseFilePath)
		if err != nil {
			log.Println(err)
			return entities.DownloadStatusFailed
		}
	}

	entities.DownloadedFilesToPathMap[downloadEntity.Id] = baseFilePath
	return entities.DownloadStatusSuccess
}

func downloadConcurrently(downloadEntity *entities.DownloadEntity) {
	downloadEntity.StartTime = time.Now()

	baseFilePath, err := createAndGetBaseFilePath(downloadEntity)
	if err != nil {
		log.Println(err)
		downloadEntity.EndTime = time.Now()
		downloadEntity.Status = entities.DownloadStatusFailed
		return
	}

	noOfFiles := len(downloadEntity.Files)

	workQueue := make(chan string, noOfFiles)
	for _, url := range downloadEntity.Files {
		workQueue <- url
	}

	parallelWorkersCount := entities.ParallelFileDownloadWorkersCount
	if parallelWorkersCount > noOfFiles {
		parallelWorkersCount = noOfFiles
	}

	for i := 0; i < parallelWorkersCount; i++ {
		go func(channel chan string) {
			for {
				url := <-channel
				err := downloadAndSaveFileLocally(url, baseFilePath)
				if err != nil {
					log.Println(err)
				}
			}
		}(workQueue)
	}
	
	downloadEntity.EndTime = time.Now()
	downloadEntity.Status = entities.DownloadStatusSuccess
}

func createAndGetBaseFilePath(downloadEntity *entities.DownloadEntity) (string, error) {
	baseFilePath := entities.DownloadFilePath + downloadEntity.Id + "/"
	err := os.Mkdir(baseFilePath, 0744) // Only the owner can read, write & execute. Everyone else can only read.
	if err != nil {
		log.Println(err)
		return "", err
	}

	return baseFilePath, nil
}



func downloadAndSaveFileLocally(url string, baseFilePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	extension := filepath.Ext(url)
	fileName := uuid.New().String() + extension
	filePath := baseFilePath + fileName

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetDownloadStatus(downloadId string) (*entities.DownloadEntity, error) {
	downloadEntity, keyExists := entities.DownloadIdToDownloadStatusMap[downloadId]

	if !keyExists {
		return nil, errors.New("unknown download ID")
	}

	return downloadEntity, nil
}
