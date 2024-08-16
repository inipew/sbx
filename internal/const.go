package internal

const (
	WorkDir        = "/etc/wth"
	BinDir          = WorkDir + "/bin"
	BackupDir       = WorkDir + "/backup"
	TmpDir          = WorkDir + "/tmp"
	LogDir          = WorkDir + "/log"
	DomainFilePath  = WorkDir + "/domain.txt"
	SingboxDir      = WorkDir + "/sing-box"
	SingboxConfDir	= SingboxDir + "/config"
	SingboxBinPath	= "/usr/local/bin/sing-box"
	SingboxLogFilePath	= "/usr/local/sing-box/sing-box.log"
	CaddyBinPath	= "/usr/bin/caddy"
	CaddyDir        = WorkDir + "/caddy"
	CaddyFilePath   = CaddyDir + "/Caddyfile"
	CaddylogFilePath = "/var/log/caddy/access.log"
	CaddyServicePath = "/etc/systemd/system/caddy.service"
	SingBoxServicePath = "/etc/systemd/system/sing-box.service"
)
const CaddyServiceContent = `[Unit]
Description=Caddy
Documentation=https://caddyserver.com/docs/
After=network.target network-online.target
Requires=network-online.target

[Service]
Type=notify
User=caddy
Group=caddy
ExecStart=/usr/bin/caddy run --environ --config /etc/caddy/Caddyfile
ExecReload=/usr/bin/caddy reload --config /etc/caddy/Caddyfile --force
TimeoutStopSec=5s
LimitNOFILE=1048576
PrivateTmp=true
ProtectSystem=full
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
`

const CaddyFileContent = `
{
        log Utama {
                output file /var/log/caddy/caddy.log {
                        roll_keep 15
                        roll_keep_for 48h
                }
                format console
        }
        servers {
                trusted_proxies static 173.245.48.0/20 103.21.244.0/22 103.22.200.0/22 103.31.4.0/22 141.101.64.0/18 108.162.192.0/18 190.93.240.0/20 188.114.96.0/20 197.234.240.0/22 198.41.128.0/17 162.158.0.0/15 104.16.0.0/13 104.24.0.0/14 172.64.0.0/13 131.0.72.0/22 2400:cb00::/32 2606:4700::/32 2803:f800::/32 2405:b500::/32 2405:8100::/32 2a06:98c0::/29 2c0f:f248::/32
        }
}

xv.akupew.toys, http://xv.akupew.toys {
        @websocket {
                header Connection *Upgrade*
                header Upgrade websocket
                header Sec-WebSocket-Key *
        }
        @http_upgrade {
                header Connection *Upgrade*
                header Upgrade websocket
                not header Sec-WebSocket-Key *
        }
        @grpc {
                header Content-Type "application/grpc"
                protocol grpc
        }

        # Mapping untuk WebSocket dan http-upgrade
        map {path} {backend} {
                /trojan 127.0.0.1:8003
                /vmess 127.0.0.1:8002
                /vless 127.0.0.1:8001
                /trojan/Tun 127.0.0.1:8007
                /vmess/Tun 127.0.0.1:8008
                /vless/Tun 127.0.0.1:8009
        }

        map {path} {backend_http_upgrade} {
                /trojan 127.0.0.1:8006
                /vmess 127.0.0.1:8005
                /vless 127.0.0.1:8004
        }

        handle @websocket {
                @rewrite_path_websocket {
                        path_regexp ^/.*?/(trojan|vmess|vless)
                }
                handle @rewrite_path_websocket {
                        rewrite * /{http.regexp.1}
                }
                reverse_proxy {backend}
        }

        handle @http_upgrade {
                @rewrite_path_http_upgrade {
                        path_regexp ^/.*?/(trojan|vmess|vless)
                }
                handle @rewrite_path_http_upgrade {
                        rewrite * /{http.regexp.1}
                }
                reverse_proxy {backend_http_upgrade}
        }
        handle @grpc {
                reverse_proxy {backend} {
                        transport http {
                                versions h2c
                        }
                }
        }

        header {
                Cache-Control "max-age=86400"
                X-Content-Type-Options "nosniff"
                X-Frame-Options "DENY"
                X-XSS-Protection "1; mode=block"
        }

        log Site {
                output file /var/log/caddy/access.log {
                        roll_keep 15
                        roll_keep_for 48h
                }
                format json
        }
}
:9090 {
        reverse_proxy 127.0.0.1:9091
}
`
const SingBoxServiceContent = `
[Unit]
Description=sing-box service
Documentation=https://sing-box.sagernet.org
After=network.target nss-lookup.target network-online.target

[Service]
CapabilityBoundingSet=CAP_NET_ADMIN CAP_NET_BIND_SERVICE CAP_SYS_PTRACE CAP_DAC_READ_SEARCH
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE CAP_SYS_PTRACE CAP_DAC_READ_SEARCH
ExecStart=/usr/local/bin/sing-box -D /var/lib/sing-box -C /usr/local/etc/sing-box run
ExecReload=/bin/kill -HUP $MAINPID
Restart=on-failure
RestartSec=10s
LimitNOFILE=infinity

[Install]
WantedBy=multi-user.target`