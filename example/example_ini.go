package main

import (
	"encoding/json"
	"fmt"
	"github.com/jacktea/go-eureka-client/eureka"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func main() {
	f := kvs.GetCurrentFilePath("app.ini", 1)
	fmt.Println(f)
	iniconf := ini.NewIniFileConfigSource(f)
	conf := kvs.NewDefaultCompositeConfigSource(iniconf)
	fmt.Println(conf.GetDefault("eureka.instance.statusPageUrlPath", ""))
	//
	client := eureka.NewClient(conf)

	data, _ := json.Marshal(client.InstanceInfo)
	fmt.Println(string(data))
	client.Start()
	c := make(chan int, 1)
	x := <-c
	fmt.Println(x)

}

//
//// disk usage of path/disk
//func DiskUsage(path string) (uint64, uint64, uint64) {
//	fs := syscall.Statfs_t{}
//	err := syscall.Statfs(path, &fs)
//	if err != nil {
//		return 0, 0, 0
//	}
//	total := fs.Blocks * uint64(fs.Bsize)
//	free := fs.Bfree * uint64(fs.Bsize)
//	used := total - free
//	return total, free, used
//}
