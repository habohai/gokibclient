package services

type UserRequest struct {
	UID    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}
