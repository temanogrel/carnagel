job "hashi-ui" {
  region = "global"
  datacenters = ["hetzner1"]
  type = "service"

  constraint {
    attribute = "${meta.special_app}"
    operator = "set_contains"
    value     = "hashi-ui"
  }

  update {
    stagger = "10s"
    max_parallel = 1
  }

  group "hashi-ui" {
    count = 1
    
    restart {
      attempts = 6
      interval = "1m"
      delay = "10s"
      mode = "delay"
    }

  }
}
