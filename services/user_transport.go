package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func GetUserInfoRequest(_ context.Context, request *http.Request, r interface{}) error {
	userRequest := r.(UserRequest)
	request.URL.Path += "/user/" + strconv.Itoa(userRequest.UID)
	return nil
}

func GetUserInfoResponse(_ context.Context, resp *http.Response) (response interface{}, err error) {
	if resp.StatusCode > 400 {
		return nil, errors.New("no data")
	}

	userResponse := UserResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, err
	}

	return userResponse, nil
}
