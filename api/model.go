package api

type response struct {
	Status string `json:"status"`
}

type uploadResult struct {
	response
	Data UploadResponse `json:"response"`
}

type UploadResponse struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}
