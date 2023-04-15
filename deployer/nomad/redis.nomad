job "redis" {
  type = "service"
  region = "global"
  datacenters = [
    "hetzner1"
  ]

   constraint {
    attribute = "${meta.special_app}"
    operator = "set_contains"
    value     = "rabbitmq"
  }

  update {
    stagger = "10s"
    max_parallel = 1
  }

  group "redis-server" {
    count = 6

    restart {
      attempts = 6
      interval = "1m"
      delay = "10s"
      mode = "delay"
    }

    task "master" {
      driver = "docker"

      config {
        image = "redis:3.2-alpine"

        args = [
          "--maxmemory", "80G",
          "--maxmemory-policy", "allkeys-lru"
        ]

        port_map {
          redis = 6379
        }
      }

      service {
        name = "redis"
        tags = ["redis"]
        port = "redis"
      }

      resources {
        cpu = 8000
        memory = 80000
        network {
          mbits = 10

          port "redis" {
            static = 6379
          }
        }
      }
    }

    task "exporter" {
      driver = "docker"

      config {
        image = "oliver006/redis_exporter"

        port_map {
          metrics = 9121
        }
      }

      env {
        REDIS_ADDR = "redis://${attr.unique.network.ip-address}:6379"
      }

      service {
        name = "redis-metrics"
        tags = ["redis", "metrics"]
        port = "metrics"
      }

      resources {
        cpu = 500
        memory = 500

        network {
          mbits = 100
          port "metrics" {}
        }
      }
    }
  }
}
