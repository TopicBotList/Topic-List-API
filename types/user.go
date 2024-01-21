package types

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
