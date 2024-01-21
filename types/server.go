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
