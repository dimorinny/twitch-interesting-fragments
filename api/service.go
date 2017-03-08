package api

func StartUploaderService(uploader Uploader, notifier <-chan struct{}) {
	for range notifier {
		uploader.Upload()
	}
}
