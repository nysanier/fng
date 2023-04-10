package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func InitConfigMap() {
	env := os.Getenv("configmap")
	fmt.Printf("configmap env: %v\n", env)

	f, err := os.OpenFile("/tmp/configmap.json", os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("open file fail, err: %v\n", err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("readall fail, err: %v\n", err)
	}

	fmt.Printf("configmap file: %v\n", string(buf))
}
