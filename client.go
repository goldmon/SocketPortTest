package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var ip string

func checkUDPPort(address string, timeout time.Duration) (bool, time.Duration) {
	var conn net.Conn
	var err error
	start := time.Now()
	conn, err = net.DialTimeout("udp", address, timeout)

	if err != nil {
		// fmt.Println("DialTimeout error:", err)
		return false, time.Since(start)
	}
	defer conn.Close()
	// Send a message to the server
	_, err = conn.Write([]byte("Hello, World!"))
	if err != nil {
		// fmt.Println("Write error:", err)
		return false, time.Since(start)
	}
	conn.Write([]byte("quit"))
	conn.Close()
	return true, time.Since(start)
}

// checkPort 检测端口是否开放
func checkTCPPort(address string, timeout time.Duration) (bool, time.Duration) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		// fmt.Println(err)
		return false, time.Since(start)
	} else {
		conn.Write([]byte("quit")) // 发送数据
	}
	conn.Close()
	return true, time.Since(start)
}

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

func Client_main(port string, ipO string) {
	//ip = "127.0.0.1:"
	ip := ipO + ":"
	conn, err := net.Dial("tcp", ip+port)
	if err != nil {
		fmt.Println("控制口连接错误: ", err)
		return
	}
	defer conn.Close() // 关闭TCP连接
	//arra := [2]string{"tcp", "udp"}
	arra := [2]string{"udp"}
	portOpen := false
	for _, value := range arra {
		fmt.Print("协议：")
		fmt.Print(value)
		fmt.Println("")
		for i := 10; i <= 65535; i++ {
			fmt.Printf("\r %d", i)
			if i == 80 || i == 443 || i == 3306 || i == 9090 || i == 61022 {
				continue
			}
			//conn.Write([]byte("close"))
			//time.Sleep(1 * time.Second)
			_, err := conn.Write([]byte("port:" + value + ":" + strconv.Itoa(i))) // 发送数据
			if err != nil {
				return
			}
			for {
				reader := bufio.NewReader(conn)
				var buf [128]byte
				n, err := reader.Read(buf[:]) // 读取数据
				if err != nil {
					fmt.Print("======>read from client failed, err: ", err)
					break
				}
				recvStr := string(buf[:n])
				if strings.HasPrefix(recvStr, "OK") { //port:tcp:1050
					portOpen = true
					break
				}
			}
			var result bool
			var usedTime time.Duration
			if portOpen {
				switch value {
				case "tcp":
					result, usedTime = checkTCPPort(ip+strconv.Itoa(i), 2*time.Second)
				case "udp":
					result, usedTime = checkUDPPort(ip+strconv.Itoa(i), 2*time.Second)
				}
				if result {
					fmt.Println("")
					fmt.Print("端口：")
					fmt.Print(i)
					fmt.Print("，结果：=====>")
					fmt.Print(result)
					fmt.Print("，时间：")
					fmt.Print(usedTime)
					fmt.Println("")
				}
			}

		}
	}
}
