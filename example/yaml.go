package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/ghodss/yaml"
)

func main() {

	file, _ := os.Getwd()
	fmt.Println("current path:", file)

	cfg := file + "/application.yml"
	data, err := readFile(cfg)

	c := make(map[string]interface{}, 0)
	err = yaml.Unmarshal([]byte(data), c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Println(c)
	v := reflect.ValueOf(c)
	kv(v, "")

}

func kv(mapv reflect.Value, parentPath string) {
	//fmt.Println(1, mapv.Kind(), mapv.Type())
	iter := mapv.MapRange()
	for iter.Next() {
		k := iter.Key()
		val := iter.Value()
		//fmt.Println(val.Kind(), val.Type(), val.Elem())
		valv := val.Elem()
		//fmt.Println(valv.Kind(), valv.Type())
		//fmt.Println(3, parentPath, k, val, valv.Kind())
		if valv.Kind() == reflect.Map {
			path := fmt.Sprintf("%s.%v", parentPath, k)
			kv(valv, path)
			continue
		} else {
			key := fmt.Sprintf("%s.%v", parentPath, k)[1:]
			value := fmt.Sprintf("%v", val)
			fmt.Println(key, value)
		}
	}

}

func readFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, err := ioutil.ReadAll(f)
	return src, err
}
