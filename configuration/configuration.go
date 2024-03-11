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
			ReturnUrl: "https://beta.topiclist.xyz/",
		},
		Client: types.Client{
			Id:       "925740376948609034",
			Secret:   os.Getenv("CLIENT_SECRET"),
			Token:    os.Getenv("CLIENT_TOKEN"),
			Callback: "https://beta.topiclist.xyz/auth/callback",
		},
		Collection: "entities",
		APIUrl:     "https://api.topiclist.xyz/",
		HCaptcha:   os.Getenv("HCAPTCHA_SECRET"),
	}
}

func GetConfig() types.Config {
	return getConfig()
}
