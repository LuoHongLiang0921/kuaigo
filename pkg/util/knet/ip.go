// @Description ip 地址

package knet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kfile"
	"net"
	"os"
	"strconv"
	"strings"
)

// GetLocalIP
//  @Description  获取本机IP
//  @Return string
//  @Return error
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("unable to determine locla ip")
}

// GetLocalMainIP
//  @Description  GetLocalMainIP
//  @Return string
//  @Return int
//  @Return error
func GetLocalMainIP() (string, int, error) {
	// UDP Connect, no handshake
	conn, err := net.Dial("udp", "8.8.8.8:8")
	if err != nil {
		return "", 0, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), localAddr.Port, nil
}

// GetMacAddrs
//  @Description  获取mac地址
//  @Return macAddrs
func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

// GetIPs
//  @Description  GetIPs
//  @Return ips
func GetIPs() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIPNet := address.(*net.IPNet)
		if isValidIPNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

// Int2IP
//  @Description  Int2IP
//  @Param nn
//  @Return net.IP
func Int2IP(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

// PrivateIP2Int
//  @Description  PrivateIP2Int
//  @Return uint32
func PrivateIP2Int() uint32 {
	ip, err := PrivateIPv4()
	if err != nil {
		return 0
	}
	return IP2Int(ip)
}

// Lower16BitPrivateIP
//  @Description  Lower16BitPrivateIP
//  @Return uint16
//  @Return error
func Lower16BitPrivateIP() (uint16, error) {
	ip, err := PrivateIPv4()
	if err != nil {
		return 0, err
	}
	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

// LocalIP
//  @Description  获取本地IP 优先使用etch0
//  @Param optionalIName 网卡命名 不传 etch0
//  @Return string
func LocalIP(optionalIName ...string) string {
	if ip := os.Getenv("BINDHOSTIP"); ip != "" {
		return ip
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
		}
	}
	var name string
	if len(optionalIName) != 0 && optionalIName[0] != "" {
		name = optionalIName[0]
	} else if name = os.Getenv("INAME"); name == "" {
		name = "eth0"
	}
	n, err := net.InterfaceByName(name)
	if err != nil {
		return ""
	}
	addrs, err := n.Addrs()
	if err != nil {
		return ""
	}
	for i := range addrs {
		if ipnet, ok := addrs[i].(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil {
				return ip.String()
			}
		}
	}
	return ""
}
// StringIpToInt
//  @Description  字符串ip转整型
//  @Param ipstring
//  @Return int
func StringIpToInt(ipstring string) int {
	ipSegs := strings.Split(ipstring, ".")
	var ipInt = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}
// IpIntToString
//  @Description  整型ip转字符串
//  @Param ipInt
//  @Return string
func IpIntToString(ipInt int) string {
	ipSegs := make([]string, 4)
	var size = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < size; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[size-i-1] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < size; i++ {
		buffer.WriteString(ipSegs[i])
		if i < size-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

// Hostname
//  @Description  获取Hostname
//  @Return string
func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return kfile.ReadFile("/etc/hostname")
	}
	return name
}

// PrivateIPv4
//  @Description  PrivateIPv4
//  @Return net.IP
//  @Return error
func PrivateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		ip := ipnet.IP.To4()
		if IsPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

// IsPrivateIPv4
//  @Description  是否为IPV4
//  @Param ip
//  @Return bool
func IsPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

// IP2Int
//  @Description  IP2Int
//  @Param ip
//  @Return uint32
func IP2Int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}
