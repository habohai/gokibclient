package util

import (
	"context"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/haibeihabo/gokibclient/services"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"

	httptransport "github.com/go-kit/kit/transport/http"

	consulapi "github.com/hashicorp/consul/api"
)

func GetUser() (string, error) {
	// 第一步： 创建client
	config := consulapi.DefaultConfig()
	config.Address = "192.168.31.82:8500" // 注册中心地址

	apiClient, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}
	client := consul.NewClient(apiClient)

	logger := log.NewLogfmtLogger(os.Stdout)

	tags := []string{"primary", "v1"}
	// 可实时查询服务实例的状态信息
	instancer := consul.NewInstancer(client, logger, "userservice", tags, true)

	factory := func(serviceURL string) (endpoint.Endpoint, io.Closer, error) {
		tat, _ := url.Parse("http://" + serviceURL)
		return httptransport.NewClient("GET", tat, services.GetUserInfoRequest, services.GetUserInfoResponse).Endpoint(), nil, nil
	}
	endpointer := sd.NewEndpointer(instancer, factory, logger)

	//mylb := lb.NewRoundRobin(endpointer) // 轮询
	mylb := lb.NewRandom(endpointer, time.Now().UnixNano()) // 随机

	getUserInfo, err := mylb.Endpoint()
	if err != nil {
		return "", err
	}

	ctx := context.Background() // 第三步：创建一个context上下文对象

	// 第四部：执行
	res, err := getUserInfo(ctx, services.UserRequest{UID: 101})
	if err != nil {
		return "", err
	}

	// 第五步：断言，得到响应值
	userInfo := res.(services.UserResponse)
	return userInfo.Result, nil
}
