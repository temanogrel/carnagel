job "rabbitmq" {
  region = "global"
  datacenters = ["hetzner1"]
  type = "service"

  constraint {
    attribute = "${meta.special_app}"
    operator = "set_contains"
    value     = "rabbitmq"
  }

  update {
    stagger = "10s"
    max_parallel = 1
  }

  group "rabbitmq" {
    count = 1
    
    restart {
      attempts = 6
      interval = "1m"
      delay = "10s"
      mode = "delay"
    }

    task "rabbitmq" {
      driver = "docker"

      config {
        image = "rabbitmq:3.6-management-alpine"
        port_map {
          clients1 = 5671
          clients2 = 5672
          management = 15672
        }

        volumes = [
          "rabbitmq-data:/var/lib/rabbitmq"
        ]

        volume_driver = "local"
      }

      env {
        RABBITMQ_DEFAULT_USER = "carnagel"
        RABBITMQ_DEFAULT_PASS = "resile-vestry-gammer-nitrate-footrest-pequod-teak-slain"
        RABBITMQ_DEFAULT_VHOST = "carnagel"
      }

      service {
        name = "rabbitmq"
        tags = ["rabbitmq"]
        port = "clients2"
      }

      service {
        name = "rabbitmq-management"
        tags = ["rabbitmq", "http"]
        port = "management"

        check {
          type     = "http"
          path     = "/"
          port     = "management"
          interval = "10s"
          timeout  = "2s"
        }
      }

      resources {
        cpu = 25000
        memory = 80000

        network {
          mbits = 1000

          port "clients1" {
            static = 5671
          }
          port "clients2" {
            static = 5672
          }
          port "management" {
            static = 15672
          }
        }
      }
    }
  }
}
