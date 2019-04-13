package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jacktea/go-eureka-client/eureka/config"

	"github.com/ghodss/yaml"
)

func main_test() {

	file, _ := os.Getwd()
	fmt.Println("current path:", file)

	cfg := file + "/application.yml"
	data, err := ReadFile(cfg)

	c := config.EurekaConfig{
		Eureka: config.Eureka{
			Client:   config.NewEurekaClientConfig(),
			Instance: config.NewEurekaInstanceConfig(),
		},
	}
	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(c.Eureka.Client.ServiceUrl)
	fmt.Println(c.Eureka.Instance.LeaseRenewalIntervalInSeconds)
	fmt.Println(c.Application.Name)
	fmt.Println(c.Application.Port)
	fmt.Println(c.Eureka.Instance.Appname)
	//	c.Eureka.Client.ServiceUrl = make(map[string]string)
	//	c.Eureka.Client.ServiceUrl["ali-hz-1"] = "http://127.0.0.1:8761/eureka"
	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
}

func ReadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, err := ioutil.ReadAll(f)
	return src, err
}
