package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// arr := []int{1,2,3,4,5}
	// for i := arr {
	// 	fmt.Println(i)
	// }
	data, err := os.ReadFile("./creds.txt")
	if err != nil {
		panic(err.Error)
	}
	var username, password string
	for _, str := range strings.Split(string(data), "\n") {
		kv := strings.Split(str, "=")
		if kv[0] == "AUTH_USERNAME" {
			username = kv[1]
		} else if kv[0] == "AUTH_PASSWORD" {
			password = kv[1]
		} else {
			fmt.Errorf("key %s is invalid", kv[0])

		}
	}
	fmt.Println(username, password)

	// fmt.Println(string(data))
}
