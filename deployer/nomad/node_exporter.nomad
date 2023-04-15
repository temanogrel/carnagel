job "node_exporter" {
  region = "global"
  datacenters = ["hetzner1"]
  type = "system"

  update {
    stagger = "1s"
    max_parallel = 100
  }

  group "node_exporter" {   
    restart {
      attempts = 6
      interval = "1m"
      delay = "10s"
      mode = "delay"
    }

    task "node_exporter" {
      driver = "raw_exec"

      config {
        command = "node_exporter-0.14.0.linux-amd64/node_exporter"
        args = [
          "-web.listen-address",
          "${NOMAD_ADDR_node_exporter}",
          "--collector.textfile.directory",
          "/var/node_exporter/textfile_collector"
        ]
      }

      artifact {
        source = "https://github.com/prometheus/node_exporter/releases/download/v0.14.0/node_exporter-0.14.0.linux-amd64.tar.gz"
      }

      service {
        name = "node-exporter"
        tags = ["node-exporter"]
        port = "node_exporter"

        check {
          type     = "http"
          path     = "/metrics"
          port     = "node_exporter"
          interval = "15s"
          timeout  = "5s"
        }
      }

      resources {
        cpu = 1000
        memory = 1024
        network {
          mbits = 50
          port "node_exporter" {
            static = "9100"
          }
        }
      }
    }
  }
}
