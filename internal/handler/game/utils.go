package game

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Profile(user entity.User, skin *entity.Skin, cape *entity.Cape) (gin.H, error) {
	timestampNow := time.Now().Unix()
	texturesUrl := os.Getenv("AWS_TEXTURES_URL")

	textures := gin.H{}
	if skin != nil {
		textures["SKIN"] = gin.H{
			"url": fmt.Sprintf("%s/%s", texturesUrl, skin.TextureHash),
		}
	}

	if cape != nil {
		textures["CAPE"] = gin.H{
			"url": fmt.Sprintf("%s/%s", texturesUrl, cape.TextureHash),
		}
	}

	texturesProperty := gin.H{
		"timestamp":   timestampNow,
		"profileId":   user.MinecraftCredential.UUID,
		"profileName": user.MinecraftCredential.Username,
		"textures":    textures,
	}

	texturesPropertyJsonString, err := json.Marshal(texturesProperty)
	if err != nil {
		return gin.H{}, err
	}

	texturesPropertyBase64 := base64.StdEncoding.EncodeToString(texturesPropertyJsonString)

	profile := gin.H{
		"id":   user.MinecraftCredential.UUID,
		"name": user.MinecraftCredential.Username,
		"properties": []any{
			gin.H{
				"name":  "textures",
				"value": texturesPropertyBase64,
			},
		},
	}

	return profile, nil
}
