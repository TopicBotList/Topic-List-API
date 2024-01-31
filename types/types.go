package types

type Server struct {
	Name        string `json:"name"`
	Icon        string `json:"servericon"`
	ID          string `json:"Servid" description:"The servers ID"`
	Owner       string `json:"owner"`
	Votes       int    `json:"ServVotes"`
	Category    string `json:"category"`
	Views       int    `json:"views"`
	Summary     string `json:"summary"`
	Description string `json:"descrip"`
	Invite      string `json:"InviteURL"`
}

type User struct {
	Bio         string      `json:"biography"`
	LongBio     string      `json:"longbio"`
	Avatar      string      `json:"avatar"`
	MfaEnabled  bool        `json:"mfa_enabled"`
	Token       string      `json:"token"`
	AccessToken string      `json:"access_token" bson:"access_token"`
	ID          string      `json:"id" description:"The users ID"`
	Username    string      `json:"username"`
	AppID       interface{} `json:"appId" bson:"appId"`
	Entity      interface{} `json:"entity" bson:"entity"`
	DisplayName string      `json:"display_name"`
}

type Bots struct {
	BotID         string `json:"botid"`
	Name          string `json:"botname"`
	Discriminator string `json:"Discriminator"`
	Website       string `json:"Website URL"`
	Github        string `json:"github"`
	Avatar        string `json:"botav"`
	Votes         bool   `json:"votes"`
	Shortdesc     string `json:"shortdesc"`
	Prefix        string `json:"prefix"`
	Publicity     string `json:"public"`
	Longdesc      string `json:"longdesc"`
	Support       string `json:"support"`
	OwnerAvatar   string `json:"ownerav"`
	OwnerName     string `json:"ownername"`
}
