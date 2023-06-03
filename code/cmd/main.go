package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Database struct {
	Server string
	Ports  []int
}

type Servers struct {
	Alpha detail
}

type detail struct {
	IP string
}

type Config struct {
	Database Database
	Servers  Servers
}

func main() {
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}

	var config Config
	viper.Unmarshal(&config)
	fmt.Printf("config %+v\n", config)
}
