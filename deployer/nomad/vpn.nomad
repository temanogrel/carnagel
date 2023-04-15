job "vpn" {
  region = "global"
  datacenters = ["hetzner1"]
  type = "service"

  constraint {
    attribute = "${meta.special_app}"
    operator = "set_contains"
    value     = "vpn"
  }

  update {
    stagger = "10s"
    max_parallel = 1
  }

  group "server" {
    count = 1

    restart {
      attempts = 6
      interval = "1m"
      delay = "10s"
      mode = "delay"
    }

    task "openvpn" {
      driver = "raw_exec"

      config {
        # manual setup required before this can work, see: https://github.com/kylemanna/docker-openvpn
        command = "docker"

        args    = [
          "run",
          "--name=openvpn",
          "-v",
          "ovpn-data:/etc/openvpn",
          "-t",
          "-p",
          "1194:1194/udp",
          "--cap-add=NET_ADMIN",
          "kylemanna/openvpn"
        ]
      }

      service {
        name = "openvpn"
        tags = ["vpn"]
        port = "openvpn"
      }

      resources {
        cpu = 500
        memory = 256
        network {
          mbits = 10
          port "openvpn" {
            static = 1194
          }
        }
      }
    }
  }
}
