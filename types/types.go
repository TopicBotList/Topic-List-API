package types

type Server struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	OwnerName   string `json:"ownerName"`
	OwnerID     string `json:"ownerID"`
	OwnerAvatar string `json:"ownerAvatar"`
	Votes       int    `json:"votes"`
	Category    string `json:"category"`
	Views       int    `json:"views"`
	Owner       string `json:"owner"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Invite      string `json:"invite"`
}

type User struct {
	Bio           string      `json:"biography"`
	LongBio       string      `json:"longbio"`
	Avatar        string      `json:"avatar"`
	Notifications []string    `bson:"notifications"`
	MfaEnabled    bool        `json:"mfa_enabled"`
	Badges        []string    `json:"badges"`
	Owner         bool        `json:"owner"`
	ZippyExpires  int64       `json:"zippyexpiredate"`
	Servers       []Server    `json:"servers"`
	Token         string      `json:"token"`
	AccessToken   string      `json:"access_token" bson:"access_token"`
	ID            string      `json:"id" description:"The users ID"`
	Username      string      `json:"username"`
	AppID         interface{} `json:"appId" bson:"appId"`
	Entity        interface{} `json:"entity" bson:"entity"`
	Name          string      `json:"name"`
	Password      string      `json:"password"`
	DisplayName   string      `json:"display_name"`
}

type Vote struct {
	Token  string `json:"token" bson:"token"`
	Server string `json:"server" bson:"server"`
	End    int64  `json:"end" bson:"end"`
}

type Bots struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Discriminator string   `json:"Discriminator"`
	Website       string   `json:"Website URL"`
	Github        string   `json:"github"`
	Avatar        string   `json:"avatar"`
	Votes         bool     `json:"votes"`
	Tags          []string `json:"tags"`
	Reviews       []string `json:"reviews"`
	Shortdesc     string   `json:"shortdesc"`
	Prefix        string   `json:"prefix"`
	Publicity     string   `json:"public"`
	Longdesc      string   `json:"longdesc"`
	Support       string   `json:"support"`
	OwnerAvatar   string   `json:"ownerAvatar"`
	OwnerName     string   `json:"ownername"`
	Analytics     string   `json:"analytics"`
}

type Config struct {
	ApiVersion int `json:"apiVersion"`
	Database   `json:"database"`
	Web        `json:"web"`
	Client     `json:"client"`
	Collection string `json:"collection"`
	APIUrl     string `json:"apiUrl"`
	HCaptcha   string `json:"hCaptcha"`
}
type Database struct {
	Url string `json:"url"`
}

type Web struct {
	Port      string `json:"port"`
	ReturnUrl string `json:"returnUrl"`
}

type Client struct {
	Id       string `json:"id"`
	Secret   string `json:"secret"`
	Token    string `json:"token"`
	Callback string `json:"callback"`
}
