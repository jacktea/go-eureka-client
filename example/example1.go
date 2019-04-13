package main

import (
	"encoding/xml"
	"net/http"

	"encoding/json"
	"github.com/jacktea/go-eureka-client/eureka"
	"strconv"
	"syscall"
)

func main_yaml() {
	client := eureka.CreateEurekaClientByYaml("application.yml")
	data, _ := xml.Marshal(client.InstanceInfo)
	println(string(data))
	client.Start()
	http.HandleFunc("/info", func(res http.ResponseWriter, req *http.Request) {

	})
	http.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		health := eureka.Health{Details: make(map[string]interface{})}
		health.Status = eureka.StatusUp
		type DiskSpace struct {
			eureka.HealthStatus
			Total     uint64
			Free      uint64
			Threshold uint64
		}

		t, f, _ := DiskUsage("/")
		ds := DiskSpace{Total: t, Free: f, Threshold: 10485760}

		//10485760
		health.Details["diskSpace"] = ds
		data, err := json.Marshal(health)
		if err != nil {
			data = []byte("{\"status\":\"UP\"}")
		}
		res.Write(data)
	})
	http.ListenAndServe(":"+strconv.Itoa(client.InstanceInfo.Port.Port), nil)
}

// disk usage of path/disk
func DiskUsage(path string) (uint64, uint64, uint64) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return 0, 0, 0
	}
	total := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	used := total - free
	return total, free, used
}
