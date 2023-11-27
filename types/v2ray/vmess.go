package v2ray

func VMessTransportFromByte(v byte) string {
	switch v {
	case 0x01:
		return "tcp"
	case 0x02:
		return "mkcp"
	case 0x03:
		return "websocket"
	case 0x04:
		return "http"
	case 0x05:
		return "domainsocket"
	case 0x06:
		return "quic"
	case 0x07:
		return "gun"
	case 0x08:
		return "grpc"
	default:
		return ""
	}
}
