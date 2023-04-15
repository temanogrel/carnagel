job "dashboards" {
  region = "global"
  type = "service"
  datacenters = [
    "hetzner1"
  ]

  constraint {
    attribute = "${meta.special_app}"
    operator = "set_contains"
    value = "dashboards"
  }

  update {
    stagger = "10s"
    max_parallel = 1
  }

  group "dashboards" {
    count = 1

    restart {
      attempts = 6
      interval = "1m"
      delay = "10s"
      mode = "delay"
    }

    task "grafana" {
      driver = "docker"

      config {
        image = "grafana/grafana:4.4.1"

        volumes = [
          "/var/lib/grafana:/var/lib/grafana"
        ]

        port_map {
          grafana = 3000
        }
      }

      service {
        name = "grafana"
        port = "grafana"
        tags = [
          "grafana",
          "http"
        ]
      }

      resources {
        cpu = 4000
        memory = 3000
        network {
          mbits = 1
          port "grafana" {
            static = "3000"
          }
        }
      }
    }

    task "hashi-ui" {
      driver = "raw_exec"

      config {
        command = "hashi-ui-linux-amd64"
      }

      artifact {
        source = "https://github.com/jippi/hashi-ui/releases/download/v0.14.0/hashi-ui-linux-amd64"
      }

      env {
        NOMAD_ENABLE = "1"
        NOMAD_ADDR = "http://127.0.0.1:4646"

        CONSUL_ENABLE = "1"
        CONSUL_ADDR = "127.0.0.1:8500"

        LISTEN_ADDRESS = "${NOMAD_ADDR_hashiui}"
      }

      service {
        name = "hashiui"
        port = "hashiui"
        tags = [
          "http",
          "monitoring",
          "production"
        ]

        check {
          type = "http"
          path = "/nomad"
          port = "hashiui"
          interval = "10s"
          timeout = "2s"
        }
      }

      resources {
        cpu = 500
        memory = 256

        network {
          mbits = 1

          port "hashiui" {
            static = 8000
          }
        }
      }
    }


    task "kibana" {
      driver = "docker"

      config {
        image = "docker.elastic.co/kibana/kibana:5.5.1"

        port_map {
          kibana = 5601
        }
      }

      env {
        ELASTICSEARCH_URL = "http://10g.es1.vee.bz:9200"

        XPACK_SECURITY_ENABLED = "false"
        XPACK_GRAPH_ENABLED = "false"
        XPACK_WATCHER_ENABLED = "false"
        XPACK_REPORTING_ENABLED = "false"
      }

      service {
        name = "kibana"
        port = "kibana"
        tags = [
          "kibana",
          "http"
        ]
      }

      resources {
        cpu = 4000
        memory = 3000

        network {
          mbits = 100

          port "kibana" {
            static = "5601"
          }
        }
      }
    }
  }
}
