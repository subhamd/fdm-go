package entities

import (
	"errors"
	"time"
)

type DownloadType string

func (dt DownloadType) IsValid() error {
	switch dt {
	case SerialDownload, ConcurrentDownload:
		return nil
	}
	return errors.New("unknown download type")
}

type DownloadStatus string

func (ds DownloadStatus) IsValid() error {
	switch ds {
	case DownloadStatusQueued, DownloadStatusFailed, DownloadStatusSuccess:
		return nil
	}
	return errors.New("unknown download status")
}

type PostDownloadsRequestObject struct {
	TYPE DownloadType `form:"type" json:"type" binding:"required"`
	URLS []string     `form:"urls" json:"urls" binding:"required"`
}

type DownloadEntity struct {
	Id           string         `form:"id" json:"id" binding:"required"`
	StartTime    time.Time      `form:"start_time" json:"start_time" binding:"required"`
	EndTime      time.Time      `form:"end_time" json:"end_time" binding:"required"`
	Status       DownloadStatus `form:"status" json:"status" binding:"required"`
	DownloadType DownloadType   `form:"download_type" json:"download_type" binding:"required"`
	Files        []string       `form:"files" json:"files" binding:"required"`
}

var DownloadIdToDownloadStatusMap = map[string]*DownloadEntity{}

var DownloadedFilesToPathMap = map[string]string{}
