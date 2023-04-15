job "collect_prometheus_size" {
  region = "global"
  datacenters = ["hetzner1"]
  type = "batch"

  constraint {
    attribute = "${meta.special_app}"
    operator = "set_contains"
    value     = "prometheus"
  }

  periodic {
    cron             = "7 */1 * * * *"
    prohibit_overlap = false
  }

  group "collect_prometheus_size" {
    task "collect_prometheus_size" {
      driver = "raw_exec"

      config {
        command = "/var/node_exporter/dir_stats.sh"
      }
    }
  }
}
