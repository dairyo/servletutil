package main

import (
	"fmt"

	"./droolswbutils"
)

func main() {
	fmt.Println(droolswbutils.Login(
		"http://192.168.50.51:8080/drools-wb/j_security_check",
		"admin",
		"admin123"))
}
