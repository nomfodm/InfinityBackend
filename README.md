# Infinity Backend [1.0.0]

## Эндпоинты и их примеры запросов

### POST /auth/signup

```go
type signUpRequest struct {
	Username string `json:"username" binding:"required,min=5,max=13"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=15"`
}
```

### POST /auth/signin

```go
type signInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
```

### POST /auth/logout

```go
type logoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
```

### POST /auth/refresh

```go
type refreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
```

### GET /auth/activate

```go
type activateRequest struct {
	ActivationCode string `json:"activationCode" binding:"required"`
}
```

### GET /user/me

```json
{
  "Authorization": "Bearer jwt token"
}
```

### POST /user/skin

```json
{
  "Authorization": "Bearer jwt token",
  "file": "skin.png"
}
```

### POST /user/cape

```json
{
  "Authorization": "Bearer jwt token",
  "file": "cape.png"
}
```

### GET /game/launcher

```json
{
  "Authorization": "Bearer jwt token"
}
```

### POST /game/join

```go
type joinRequest struct {
	AccessToken     string `json:"accessToken" binding:"required,len=32"`
	SelectedProfile string `json:"selectedProfile" binding:"required,len=32"`
	ServerID        string `json:"serverId" binding:"required,gte=39"`
}
```

### GET /game/profile/:uuid

### GET /game/hasJoined

```
?username=fsdfsd&serverId=fsdfsdf
```

### .env example

```
JWT_SECRET="123456"

# minutes
JWT_LIFETIME=30

# days
REFRESH_TOKEN_LIFETIME=14


FRONTEND_URL="http://localhost:5173"


AWS_ACCESS=123
AWS_SECRET=123
AWS_BUCKET_NAME="infinity"
AWS_URL="s3.storage.selcloud.ru"


AWS_CONTENT_URL="https://storage.infinityserver.ru"
AWS_TEXTURES_URL="https://storage.infinityserver.ru/textures"
```
