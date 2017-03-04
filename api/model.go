package api

type response struct {
	Status string `json:"status"`
}

type uploadResult struct {
	Url  string `json:"status"`
	Name string `json:"name"`
}

type UploadResponse struct {
	response
	Data uploadResult `json:"name"`
}
