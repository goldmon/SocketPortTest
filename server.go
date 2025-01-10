package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

var quitRoutine bool

/**
func main() {
	ports := []int{22, 23, 53, 80, 443, 10080, 3306, 25}
	addr := "192.168.140.3"
	timeout := 1 * time.Second

	for _, port := range ports {
		address := fmt.Sprintf("%s:%d", addr, port)

		// 检测TCP端口
		isOpen, duration := checkTCPPort(address, timeout)
		fmt.Printf("TCP %s open: %t (checked in %s)\n", address, isOpen, duration)
	}
	for _, port := range ports {
		address := fmt.Sprintf("%s:%d", addr, port)

		// 检测UDP端口
		isOpen, duration := checkUDPPort(address, timeout)
		fmt.Printf("UDP %s open: %t (checked in %s)\n", address, isOpen, duration)
	}
}
*/
// TCP Server端测试
// 处理函数
func process(port string, protocol string) {
	listen, err := net.Listen(protocol, "0.0.0.0:"+port)
	fmt.Println("开始监听端口：" + port)
	if err != nil {
		//fmt.Println("Listen() failed, err: ", err)
		return
	}
	for {
		conn, err := listen.Accept() // 监听客户端的连接请求
		if err != nil {
			fmt.Println("Accept() failed, err: ", err)
			continue
		}
		defer conn.Close() // 关闭连接
		//go process(conn) // 启动一个goroutine来处理客户端的连接请求
		for {
			reader := bufio.NewReader(conn)
			var buf [128]byte
			_, err := reader.Read(buf[:]) // 读取数据
			if err != nil {
				//fmt.Println("read from client failed, err: ", err)
				break
			}
			/**recvStr := string(buf[:n])
			if len(recvStr) > 0 {
				fmt.Println("protocol:" + protocol + ",Port:" + port + "    " + recvStr)
			}*/
			if quitRoutine {
				//fmt.Println("quit2")
				break
			}
			//time.Sleep(1 * time.Second)
			//fmt.Println(port)
		}
		//time.Sleep(1 * time.Second)
		//fmt.Println("hahah")
		if quitRoutine {
			//fmt.Println("quit0")
			break
		}
	}
	//fmt.Println("quit_all")
	listen.Close()
}
func Server_main(port string) {
	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Listen() failed, err: ", err)
		return
	}
	for {
		conn, err := listen.Accept() // 监听客户端的连接请求
		if err != nil {
			fmt.Println("Accept() failed, err: ", err)
			continue
		}
		quitRoutine = true
		time.Sleep(1 * time.Second)
		quitRoutine = false
		//go process(conn) // 启动一个goroutine来处理客户端的连接请求
		for {
			reader := bufio.NewReader(conn)
			var buf [128]byte
			n, err := reader.Read(buf[:]) // 读取数据
			if err != nil {
				fmt.Println("read from client failed, err: ", err)
				break
			}
			recvStr := string(buf[:n])
			switch recvStr {
			case "heartbeat":
				conn.Write([]byte(recvStr)) // 发送数据
			case "close":
				quitRoutine = true
				fmt.Println("close")
			default:
				if strings.HasPrefix(recvStr, "port:") { //port:tcp:1050
					quitRoutine = true
					time.Sleep(50 * time.Millisecond)
					quitRoutine = false
					str := strings.Split(recvStr, ":")
					fmt.Println(str[1])
					//conn.Write([]byte(str[1])) // 发送数据
					conn.Write([]byte("OK")) // 发送数据
					go process(str[2], str[1])
				}
			}

		}
	}
}
