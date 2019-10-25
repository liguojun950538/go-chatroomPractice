package common

type UserCanBeSeen struct {
	UserID     int    `json:"userid"`
	UserName   string `json:"username"`
	UserStatus int    `json:"userstatus"`
}

type UserCanNotBeSeen struct {
	Passwd string `json:"passwd"`
}
type User struct {
	UserCanBeSeen
	UserCanNotBeSeen
}
