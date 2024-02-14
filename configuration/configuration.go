package configuration

import (
	"os"

	"go.topiclist.xyz/types"
)

func getConfig() types.Config {
	return types.Config{
		ApiVersion: 5,
		Database: types.Database{
			Url: os.Getenv("DATABASE_URL"),
		},
		Web: types.Web{
			Port:      "8080",
			ReturnUrl: "http://localhost:3000",
		},
		Client: types.Client{
			Id:       "1006423342401732691",
			Secret:   os.Getenv("CLIENT_SECRET"),
			Token:    os.Getenv("CLIENT_TOKEN"),
			Callback: "http://127.0.0.1:8080//auth/callback",
		},
		Collection: "entities",
		APIUrl:     "http://127.0.0.1:8080/",
		HCaptcha:   os.Getenv("HCAPTCHA_SECRET"),
	}
}

func GetConfig() types.Config {
	return getConfig()
}
