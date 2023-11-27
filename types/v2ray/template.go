package v2ray

const (
	configTemplate = `
{
    "api": {
        "services": [
            "StatsService"
        ],
        "tag": "api"
    },
    "inbounds": [
        {
            "listen": "127.0.0.1",
            "port": {{ .API.Port }},
            "protocol": "dokodemo-door",
            "settings": {
                "address": "127.0.0.1"
            },
            "tag": "api"
        },
        {
            "listen": "127.0.0.1",
            "port": {{ .Proxy.Port }},
            "protocol": "socks",
            "settings": {
                "ip": "127.0.0.1",
                "udp": true
            },
            "sniffing": {
                "destOverride": [
                    "http",
                    "tls"
                ],
                "enabled": true
            },
            "tag": "proxy"
        }
    ],
    "log": {
        "loglevel": "none"
    },
    "outbounds": [
        {
            "protocol": "vmess",
            "settings": {
                "vnext": [
                    {
                        "address": "{{ .VMess.Address }}",
                        "port": {{ .VMess.Port }},
                        "users": [
                            {
                                "alterId": 0,
                                "id": "{{ .VMess.ID }}"
                            }
                        ]
                    }
                ]
            },
            "streamSettings": {
                "network": "{{ .VMess.Transport }}"
            },
            "tag": "vmess"
        }
    ],
    "policy": {
        "levels": {
            "0": {
                "downlinkOnly": 0,
                "uplinkOnly": 0
            }
        },
        "system": {
            "statsOutboundDownlink": true,
            "statsOutboundUplink": true
        }
    },
    "routing": {
        "rules": [
            {
                "inboundTag": [
                    "api"
                ],
                "outboundTag": "api",
                "type": "field"
            }
        ]
    },
    "stats": {},
    "transport": {
        "dsSettings": {},
        "grpcSettings": {},
        "gunSettings": {},
        "httpSettings": {},
        "kcpSettings": {},
        "quicSettings": {
            "security": "chacha20-poly1305"
        },
        "tcpSettings": {},
        "wsSettings": {}
    }
}
`
)
