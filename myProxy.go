/*
#!/usr/bin/env gorun
# _author_ lbbxsxlz@gmail.com
# connect with newofcortexm3@163.com
*/

package main

import (
	"os"
	"os/signal"
	"syscall"
	"bufio"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"io"
	"net"
	"sync"
)

const socks_v5 byte = 0x05
const CMD_CONNECT byte = 0x01
const addrTypeIPv4 byte = 0x01
const addrTypeDomain byte = 0x03
const addrTypeIpv6 byte = 0x04

/* SOCKS协议处理 */
func socks5Deal(reader *bufio.Reader, conn net.Conn)(string, error) {
	/*
	 * +----+----------+----------+
	 * |VER | NMETHODS | METHODS  |
	 * +----+----------+----------+
	 * | 1  |    1     | 1 to 255 |
	 * +----+----------+----------+
	 */

	var err error
	var ver byte
	if ver, err = reader.ReadByte(); err != nil || ver != socks_v5 {
		fmt.Println("socks version error: %v,%d", err, ver)
		return "", errors.New("Not socks5 protol!")
	}

	var nmethods byte
	if nmethods, err = reader.ReadByte(); err != nil || nmethods == 0 {
		fmt.Println("read nmethods fail, nmethods = %d!", nmethods)
		return "", errors.New("read nmethods fail!")
	}

	/*
	 * The values currently defined for METHOD are:
	 * o  X'00' NO AUTHENTICATION REQUIRED (supported)
	 * o  X'01' GSSAPI SSH支持的一种验证方式
	 * o  X'02' USERNAME/PASSWORD (supported)
	 * o  X'03' to X'7F' IANA ASSIGNED
	 * o  X'80' to X'FE' RESERVED FOR PRIVATE METHODS
	 * o  X'FF' NO ACCEPTABLE METHODS
	 */
	authMethods := make([]byte, nmethods)
	if _, err = io.ReadFull(reader, authMethods); err != nil {
		fmt.Println("get authMethods fail!")
		return "", errors.New("read detail methods!")
	}

	for authMethod := range authMethods {
		if authMethod == 0x00 {
			//log.Printf("NO AUTHENTICATION REQUIRED!")
			break
		}
	}

	/*
	 * +----+--------+
	 * |VER | METHOD |
	 * +----+--------+
	 * | 1  |   1    |
	 * +----+--------+
	 */
	/* 回复: socks5协议，认证方式为0，即无需密码访问 */
	response := []byte{5, 0}
	conn.Write(response)

	/*
	 * +----+-----+-------+------+----------+----------+
	 * |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	 * +----+-----+-------+------+----------+----------+
	 * | 1  |  1  | X'00' |  1   | Variable |    2     |
	 * +----+-----+-------+------+----------+----------+
	 */
	/* o  VER    protocol version: X'05' */
	if ver, err = reader.ReadByte(); err != nil || ver != socks_v5 {
		fmt.Println("socks version error: %v,%d", err, ver)
		return "", errors.New("Not socks5 protol!")
	}

	/*
	 * o  CMD
	 *	 o  CONNECT X'01' ( define in CMD_CONNECT )
	 *	 o  BIND X'02'    ( define in CMD_BIND )
	 *	 o  UDP ASSOCIATE X'03'  ( define in CMD_UDP )
	 */
	var cmd byte
	if cmd, err = reader.ReadByte(); err != nil {
		fmt.Println("can not get detail cmd")
		return "", errors.New("get cmd fail!")
	}

	if cmd != CMD_CONNECT {
		return "", errors.New("Only support connect!")
	}

	/*  o  RSV    RESERVED */
	/* 跳过RSV保留字段 */
	if _, err = reader.ReadByte(); err != nil {
		fmt.Println("read reserved byte fail!")
		return "", errors.New("skip RSV fail")
	}

	/*
	 * o  ATYP   address type of following address
	 *	 o  IP V4 address: X'01'
	 *	 o  DOMAINNAME: X'03'
	 *	 o  IP V6 address: X'04'
	 * o  DST.ADDR       desired destination address
	 * o  DST.PORT desired destination port in network octet
	 *	 order
	 *
	 */
	var addrType byte
	if addrType, err = reader.ReadByte(); err != nil {
		fmt.Println("can not get address type!")
		return "", errors.New("get addr type fail")
	}

	//log.Printf("客户端请求的远程服务器地址类型是:%d", addrtype)
	var strAddr string
	switch addrType {
	case addrTypeIPv4:
		ipv4 := make([]byte, 4)
		if _, err = io.ReadFull(reader, ipv4); err != nil {
			fmt.Println("can not get ipv4 address!")
			return "", errors.New("get Ipv4 addr fail")
		}

		strAddr = net.IP(ipv4).String()

	case addrTypeDomain:
		var nameLen byte
		if nameLen, err = reader.ReadByte(); err != nil {
			fmt.Println("can not get domain name's length!")
			return "", errors.New("get domain name length fail")
		}

		domainName := make([]byte, nameLen)
		if _, err = io.ReadFull(reader, domainName); err != nil {
			fmt.Println("can not get domain name!")
			return "", errors.New("get domain name fail")
		}

		strAddr = string(domainName)
	default:
		fmt.Println("Do not deal with the IPv6 address for now!")
		return "", errors.New("Not support IPv6 addr!")
	}

	var port int16
	if err = binary.Read(reader, binary.BigEndian, &port); err != nil {
		fmt.Println("can not get the prot")
		return "", errors.New("get port fail")
	}

	/*
	 * +----+-----+-------+------+----------+----------+
	 * |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	 * +----+-----+-------+------+----------+----------+
	 * | 1  |  1  | X'00' |  1   | Variable |    2     |
	 * +----+-----+-------+------+----------+----------+
	 *
	 * Where:
	 * o  VER    protocol version: X'05'
	 * o  REP    Reply field:
	 *	 o  X'00' succeeded
	 *	 o  X'01' general SOCKS server failure
	 *	 o  X'02' connection not allowed by ruleset
	 *	 o  X'03' Network unreachable
	 *	 o  X'04' Host unreachable
	 *	 o  X'05' Connection refused
	 * 	 o  X'06' TTL expired
	 *	 o  X'07' Command not supported
	 *	 o  X'08' Address type not supported
	 *	 o  X'09' to X'FF' unassigned
	 * o  RSV    RESERVED (must be set to 0x00)
	 * o  ATYP   address type of following address
	 * o  IP V4 address: X'01'
	 *	 o  DOMAINNAME: X'03'
	 *	 o  IP V6 address: X'04'
	 * o  BND.ADDR       server bound address
	 * o  BND.PORT       server bound port in network octet order
	 *
	 */
	//返回 IPv4，即ATYP=0x01,共10个字节
	response = []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	conn.Write(response)

	return fmt.Sprintf("%s:%d", strAddr, port), nil
}

func proxyHandle(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	/* socks5协议解析 */
	addr, err := socks5Deal(reader, conn)
	if err != nil {
		log.Print(err)
	}
	//log.Print("远程连接地址：", addr)

	//尝试建立远程连接
	var remote net.Conn

	remote, err = net.Dial("tcp", addr)
	if err != nil {
		log.Print(err)
		return
	}

	defer remote.Close()

	//等待组进行2个任务的同步
	wGroup := new(sync.WaitGroup)
	wGroup.Add(2)

	//goroutine
	go func() {
		defer wGroup.Done()
		//获取原地址请求，发送给目标主机
		io.Copy(remote, reader)
	}()

	go func() {
		defer wGroup.Done()
		//目标主机的数据转发给客户机
		io.Copy(conn, remote)
	}()
	wGroup.Wait()
}

func signalInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGKILL)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("receive signal to exit!")
		os.Exit(0)
	}()
}

func main() {
	var listen string

	signalInterrupt()
	flag.StringVar(&listen, "listen", "0.0.0.0:1080", "listen address(host:port)")
	flag.Parse()
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go proxyHandle(conn)
	}
}
