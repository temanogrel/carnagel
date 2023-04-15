job "minerva" {
	region = "global"
	datacenters = ["hetzner1"]
	type = "service"

	constraint {
    attribute = "${meta.special_app}"
    operator = "set_contains"
    value     = "minerva"
  }

	update {
		stagger = "10s"
		max_parallel = 1
	}

	group "minerva" {
		count = 1
		
		restart {
			attempts = 6
			interval = "1m"
			delay = "10s"
			mode = "delay"
		}

		task "minerva" {
			driver = "raw_exec"

			config {
				command = "minerva"
			}

			artifact {
        source = "http://deployer:eZKpb4uR9cEW9A7tvQF2ug3LCNZxxpzysuDYYBzbUn@teamcity.misc.vee.bz/repository/download/Minerva_Build/5738:id/minerva.zip"
      }

			env {
        DB_USER = "minerva"
        DB_PASS = "minerva"
        DB_NAME = "minerva"
        DB_PORT = "5432"
        DB_HOST = "10g.postgres.vee.bz"
        PROMETHEUS_URI = "http://10.0.30.9:29582"
        MIGRATIONS_PATH = "./"
    	}

			service {
				name = "minerva"
				tags = ["minerva"]
				port = "grpc"
			}

			resources {
				cpu = 500
				memory = 256
				network {
					mbits = 1
					port "grpc" {}
				}
			}
		}
	}
}
