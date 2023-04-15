job "minion" {
  region = "global"
  datacenters = ["hetzner1"]
  type = "system"

  constraint {
    attribute = "${meta.special_app}"
    operator = "set_contains"
    value     = "minion"
  }

  update {
    stagger = "5s"
    max_parallel = 25
  }

  group "minion" {
    restart {
      attempts = 6
      interval = "1m"
      delay = "10s"
      mode = "delay"
    }

    task "minion" {
      driver = "raw_exec"

      config {
        command = "minion"
      }

      artifact {
        source = "http://deployer:eZKpb4uR9cEW9A7tvQF2ug3LCNZxxpzysuDYYBzbUn@teamcity.misc.vee.bz/repository/download/Minion_Build/5839:id/minion.zip"
      }

      service {
        name = "minion"
        tags = ["minion"]
        port = "http"
      }

      resources {
        cpu = 500
        memory = 256
        network {
          mbits = 1
          port "http" {
            static = "6000"
          }
        }
      }
    }
  }
}
