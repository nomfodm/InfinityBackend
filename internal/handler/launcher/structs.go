package launcher

type registerUpdateRequest struct {
	Version     string `json:"version" binding:"required"`
	DownloadUrl string `json:"downloadUrl" binding:"required"`
	SHA256      string `json:"sha256" binding:"required"`
	Mandatory   bool   `json:"mandatory" binding:"required"`
}
