job "metrics" {
  region = "global"
  datacenters = ["hetzner1"]
  type = "service"

  group "mysql" {
    count = 1

    constraint {
      attribute = "${meta.special_app}"
      operator = "set_contains"
      value     = "mysql"
    }

    task "mysql-exporter" {
      driver = "docker"

      config {
        image = "prom/mysqld-exporter:v0.10.0"
        port_map {
          metrics = 9104
        }
      }

      env {
        DATA_SOURCE_NAME = "mysql_exporter:mysql_exporter@tcp(10g.mysql.vee.bz:3306)/"
      }

      service {
        name = "mysql-metrics"
        tags = ["mysql", "metrics"]
        port = "metrics"
      }

      resources {
        cpu = 100
        memory = 512

        network {
          mbits = 5
          port "metrics" {}
        }
      }
    }
  }

  group "postgres" {
    count = 1

    constraint {
      attribute = "${meta.special_app}"
      operator = "set_contains"
      value     = "postgres"
    }

    task "postgres-exporter" {
      driver = "docker"

      config {
        image = "wrouesnel/postgres_exporter:v0.2.0"
        port_map {
          metrics = 9187
        }
      }

      env {
        DATA_SOURCE_NAME = "postgresql://postgres_exporter:postgres_exporter@10g.postgres.vee.bz:5432/minerva?sslmode=disable"
      }

      service {
        name = "postgres-metrics"
        tags = ["postgres", "metrics"]
        port = "metrics"
      }

      resources {
        cpu = 100
        memory = 512

        network {
          mbits = 5
          port "metrics" {}
        }
      }
    }
  }

  group "rabbitmq" {
    count = 1

    constraint {
      attribute = "${meta.special_app}"
      operator = "set_contains"
      value     = "rabbitmq"
    }

    task "rabbitmq-exporter" {
      driver = "docker"

      config {
        image = "kbudde/rabbitmq-exporter"
        port_map {
          metrics = 9090
        }
      }

      env {
        RABBIT_URL = "http://10g.rabbitmq.vee.bz:15672"
        RABBIT_USER = "carnagel"
        RABBIT_PASSWORD = "resile-vestry-gammer-nitrate-footrest-pequod-teak-slain"
        PUBLISH_PORT = "9090"
      }

      service {
        name = "rabbitmq-metrics"
        tags = ["rabbitmq", "metrics"]
        port = "metrics"
      }

      resources {
        cpu = 100
        memory = 256

        network {
          mbits = 25
          port "metrics" {}
        }
      }
    }
  }
}
