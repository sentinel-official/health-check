package utils

import (
	"net"
)

func GetFreeTCPPort() (uint16, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return 0, err
	}

	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer conn.Close()

	return uint16(conn.Addr().(*net.TCPAddr).Port), nil
}
