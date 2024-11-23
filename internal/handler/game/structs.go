package game

type joinRequest struct {
	AccessToken     string `json:"accessToken" binding:"required,len=32"`
	SelectedProfile string `json:"selectedProfile" binding:"required,len=32"`
	ServerID        string `json:"serverId" binding:"required,gte=39"`
}
