package util

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

func BytesToUint64(v []byte) (uint64, error) {
	u, err := strconv.Atoi(string(v))
	return uint64(u), err
}

func Uint64ToBytes(u uint64) []byte {
	return []byte(fmt.Sprint(u))
}

func TxHash(b []byte) []byte {
	return crypto.Keccak256(b)
}

// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
// TODO there must be a better way to get external IP
func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if skipInterface(iface) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip := addrToIP(addr)
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func skipInterface(iface net.Interface) bool {
	if iface.Flags&net.FlagUp == 0 {
		return true // interface down
	}
	if iface.Flags&net.FlagLoopback != 0 {
		return true // loopback interface
	}
	return false
}

func addrToIP(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip
}

func SetEnvToViper(vp *viper.Viper, key string) []string {
	v := os.Getenv(key)
	var keys []string
	pairs := strings.Split(v, ",")
	for _, pair := range pairs {
		pair := strings.TrimSpace(pair)
		if len(pair) == 0 {
			break
		}
		kv := strings.Split(pair, "=")
		vp.Set(kv[0], kv[1])
		keys = append(keys, kv[0])
	}
	return keys
}
