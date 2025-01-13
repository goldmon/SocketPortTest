package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
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
func processTCP(port string) {
	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	//fmt.Println("开始监听端口：" + port)
	if err != nil {
		//fmt.Println("Listen() failed, err: ", err)
		return
	}
	defer listen.Close()
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
			n, err := reader.Read(buf[:]) // 读取数据
			if err != nil {
				//fmt.Println("read from client failed, err: ", err)
				break
			}
			recvStr := string(buf[:n])
			if recvStr == "quit" {
				//fmt.Println("read from client failed, err: ", err)
				quitRoutine = true
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
		conn.Close()
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
func closePort(port string, protocol string) {
	p1, err := strconv.Atoi(port)
	p2 := p1 - 1
	p3 := strconv.Itoa(p2)
	conn, err := net.DialTimeout(protocol, "127.0.0.1:"+p3, 2*time.Second)
	if err == nil {
		conn.Write([]byte("quit")) // 发送数据
		conn.Close()
	}
}
func processUDP(port string) {
	porti, _ := strconv.Atoi(port)
	// create udp server
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: porti,
	})
	if err != nil {
		fmt.Printf("listem failed, err: %v\n", err)
		return
	}

	defer listen.Close()
	for {
		var buf [1024]byte
		n, _, err := listen.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Printf("read data failed, err: %v\n", err)
			return
		}
		recvStr := string(buf[:n])

		if recvStr == "quit" {
			//fmt.Println("read from client failed, err: ", err)
			quitRoutine = true
			break
		} else {
			fmt.Println("客户端: ", recvStr)
		}
		/**recvStr := string(buf[:n])
		if len(recvStr) > 0 {
			fmt.Println("protocol:" + protocol + ",Port:" + port + "    " + recvStr)
		}*/
		if quitRoutine {
			//fmt.Println("quit2")
			break
		}
	}
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
					//fmt.Println(str[1])
					//conn.Write([]byte(str[1])) // 发送数据
					conn.Write([]byte("OK")) // 发送数据
					if str[1] == "tcp" {
						go processTCP(str[2])
					} else {
						go processUDP(str[2])
					}
					go closePort(str[2], str[1])
				}
			}

		}
	}
}
