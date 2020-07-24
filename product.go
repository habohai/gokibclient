package main

import (
	"fmt"
	"gomicro2/util"
	"log"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func main() {
	configA := hystrix.CommandConfig{
		Timeout:                2000,
		MaxConcurrentRequests:  5,
		RequestVolumeThreshold: 3,
		ErrorPercentThreshold:  20,
		SleepWindow:            int(time.Second * 100),
	}

	hystrix.ConfigureCommand("getuser", configA)
	err := hystrix.Do("getuser", func() error {
		res, err := util.GetUser()
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	}, func(err error) error {
		fmt.Println("降级用户")
		return err
	})

	if err != nil {
		log.Fatal(err)
	}
}
