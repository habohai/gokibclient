package main

import (
	"context"
	"fmt"
	"gomicro2/services"
	"net/url"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	target, _ := url.Parse("http://localhost:9050")
	client := httptransport.NewClient("GET", target, services.GetUserInfoRequest, services.GetUserInfoResponse)
	getUserInfo := client.Endpoint()

	ctx := context.Background()
	res, err := getUserInfo(ctx, services.UserRequest{UID: 101})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	userInfo := res.(services.UserResponse)
	fmt.Println(userInfo.Result)
}
