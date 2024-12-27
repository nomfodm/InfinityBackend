package launcher

type updateRequest struct {
	ClientVersion string `json:"clientVersion" binding:"required"`
	ClientHash    string `json:"clientHash" binding:"required"`
}
