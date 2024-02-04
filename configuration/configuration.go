package configuration

import (
	"go.topiclist.xyz/types"
)

func getConfig() types.Config {
	return types.Config{
		ApiVersion: 5,
		Database: types.Database{
			Url: "mongodb+srv://Admin:RanveerSoni11@topic.q8qcpfz.mongodb.net",
		},
		Web: types.Web{
			Port:      "8080",
			ReturnUrl: "https://servers.topiclist.xyz",
		},
		Client: types.Client{
			Id:       "1006423342401732691",
			Secret:   "rn06jIkC0J9XCuuejaRNJ3ZHGWon3gRQ",
			Token:    "MTAwNjQyMzM0MjQwMTczMjY5MQ.GYtWnK.bnr1UB0LJ-fBphrzBHL3WEFnRaQhMTx_mdBY6M",
			Callback: "https://servers.topiclist.xyz/private/auth/callback",
		},
		Collection: "entities",
		APIUrl:     "https://k02hrtapiv5j.topiclist.xyz/private",
		HCaptcha:   "be91bfbf-20c4-47de-a3af-846d5156d39c",
	}
}

func GetConfig() types.Config {
	return getConfig()
}
