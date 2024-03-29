package utils

import (
	"fmt"
	"net"
	"strings"
)

// GetOutboundIP 获取本地对外IP
func GetOutboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	loacalAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(loacalAddr.String())
	ip = strings.Split(loacalAddr.IP.String(), ":")[0]
	return
}
