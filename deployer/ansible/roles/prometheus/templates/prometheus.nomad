job "prometheus" {
  region = "global"
  datacenters = ["hetzner1"]
  type = "service"

  constraint {
    attribute = "${node.unique.id}"
    operator = "set_contains"
    value = "64912e8b-d318-d562-7a91-d115e7b4d7bb"
  }

  update {
    stagger = "10s"
    max_parallel = 1
  }

  group "prometheus" {
    count = 1

    restart {
      attempts = 6
      interval = "1m"
      delay = "10s"
      mode = "delay"
    }

    task "alertmanager" {
      driver = "raw_exec"

      config {
        command = "alertmanager-0.7.1.linux-amd64/alertmanager"
        args = [
          "-config.file",
          "/etc/alertmanager/config.yml",
          "-web.listen-address",
          "${NOMAD_ADDR_alertmanager}",
          "-storage.path",
          "/var/alertmanager"
        ]
      }

      artifact {
        source = "https://github.com/prometheus/alertmanager/releases/download/v0.7.1/alertmanager-0.7.1.linux-amd64.tar.gz"
      }

      service {
        name = "alertmanager"
        tags = ["prometheus", "alertmanager"]
        port = "alertmanager"
      }

      resources {
        cpu = 12000
        memory = 8192
        network {
          mbits = 100
          port "alertmanager" {
            static = 9091
          }
        }
      }
    }

    task "prometheus" {
      driver = "raw_exec"

      config {
        command = "prometheus-1.7.1.linux-amd64/prometheus"
        args = [
          "-web.listen-address",
          "${NOMAD_ADDR_prometheus}",
          "-config.file",
          "/etc/prometheus/prometheus.yml",
          "-storage.local.path",
          "/var/prometheus",
          "-storage.local.retention",
          "300h",
          "-storage.local.target-heap-size",
          "40187091200",
          "-alertmanager.url",
          "http://10.0.30.9:9091"
        ]
      }

      artifact {
        source = "https://github.com/prometheus/prometheus/releases/download/v1.7.1/prometheus-1.7.1.linux-amd64.tar.gz"
      }

      service {
        name = "prometheus"
        tags = ["prometheus"]
        port = "prometheus"
      }

      resources {
        cpu = 12000
        memory = 8192
        network {
          mbits = 100
          port "prometheus" {
            static = 9090
          }
        }
      }
    }
  }
}
