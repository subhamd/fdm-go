package entities

const (
	SerialDownload        DownloadType   = "serial"
	ConcurrentDownload    DownloadType   = "concurrent"
	DownloadStatusQueued  DownloadStatus = "QUEUED"
	DownloadStatusFailed  DownloadStatus = "FAILED"
	DownloadStatusSuccess DownloadStatus = "SUCCESS"

	DownloadFilePath string = "/Users/subhamd/Desktop/downloaded_file/"
	ParallelFileDownloadWorkersCount = 10
)