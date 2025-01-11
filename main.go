package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "server":
			if len(os.Args) > 2 {
				Server_main(os.Args[2])
			} else {
				fmt.Println("缺少参数")
			}
		case "client":
			if len(os.Args) > 3 {
				Client_main(os.Args[2], os.Args[3])
			} else {
				fmt.Println("缺少参数")
			}
		default:
			fmt.Println("参数不对")
		}
	} else {
		fmt.Println("请输入参数:server/client 端口 IP(服务端忽略)")
	}
}

