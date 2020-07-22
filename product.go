package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"

	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	consulapi "github.com/hashicorp/consul/api"

	"gomicro2/services"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	// 第一步： 创建client

	config := consulapi.DefaultConfig()
	config.Address = "192.168.31.82:8500" // 注册中心地址

	apiClient, _ := consulapi.NewClient(config)
	client := consul.NewClient(apiClient)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
	}
	{
		tags := []string{"primary", "v1"}
		// 可实时查询服务实例的状态信息
		instancer := consul.NewInstancer(client, logger, "userservice", tags, true)

		factory := func(serviceURL string) (endpoint.Endpoint, io.Closer, error) {
			tat, _ := url.Parse("http://" + serviceURL)
			return httptransport.NewClient("GET", tat, services.GetUserInfoRequest, services.GetUserInfoResponse).Endpoint(), nil, nil
		}
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		endpionts, _ := endpointer.Endpoints()
		fmt.Println("服务有：", len(endpionts), "个")
		if len(endpionts) == 0 {
			return
		}

		//mylb := lb.NewRoundRobin(endpointer)
		mylb := lb.NewRandom(endpointer, time.Now().UnixNano())

		for {
			// getUserInfo := endpionts[0] // 写死获取第一个
			getUserInfo, _ := mylb.Endpoint()
			ctx := context.Background()
			res, err := getUserInfo(ctx, services.UserRequest{UID: 101})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			userInfo := res.(services.UserResponse)
			fmt.Println(userInfo.Result)
			time.Sleep(time.Second * 3)
		}
	}
}
