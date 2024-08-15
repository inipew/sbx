package caddyfile

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const caddyfilePath = "/etc/caddy/Caddyfile"

func ApplyCaddyfile() error {
    cmd := exec.Command("/usr/bin/caddy", "reload", "--config", caddyfilePath, "--force")
    out, err := cmd.CombinedOutput() // Mengambil output error juga
    if err != nil {
        return fmt.Errorf("gagal memuat ulang Caddyfile: %s", string(out))
    }
    return nil
}

func ReplaceDomain(oldDomain, newDomain string) error {
    // Baca isi Caddyfile
    data, err := os.ReadFile(caddyfilePath)
    if err != nil {
        return fmt.Errorf("gagal membaca Caddyfile: %w", err)
    }

    // Ganti domain lama dengan domain baru
    updatedData := strings.ReplaceAll(string(data), oldDomain, newDomain)

    // Tulis kembali ke Caddyfile
    err = os.WriteFile(caddyfilePath, []byte(updatedData), 0644)
    if err != nil {
        return fmt.Errorf("gagal menulis Caddyfile: %w", err)
    }

    return ApplyCaddyfile()
}
const DefaultCaddyfileConfig = `
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

domain.com, http://domain.com {
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
`

// InstallDefaultConfig menulis konfigurasi default ke Caddyfile
func InstallDefaultConfig() error {
    err := os.WriteFile(caddyfilePath, []byte(DefaultCaddyfileConfig), 0644)
    if err != nil {
        return fmt.Errorf("gagal menulis Caddyfile default: %w", err)
    }
    return nil
}
