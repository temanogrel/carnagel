job "btc-plutus" {
  type = "service"
  region = "global"
  datacenters = [
    "hetzner1"
  ]

  constraint {
    attribute = "${meta.special_app}"
    operator = "set_contains"
    value = "btc"
  }

  group "plutus" {
    task "plutus" {
      driver = "raw_exec"

      config {
        command = "plutus"
      }

      artifact {
        source = "http://deployer:eZKpb4uR9cEW9A7tvQF2ug3LCNZxxpzysuDYYBzbUn@teamcity.misc.vee.bz/repository/download/Bitcoin_Plutus_Build/{{ build_number }}/plutus.zip"
      }

      env {
        WALLET_RPC_CRT = "{{ lookup('env', 'BTC_WALLET_RPC_CRT') }}"
      }

      service {
        name = "plutus"
        port = "grpc"
        tags = [
          "plutus"
        ]
      }

      resources {
        cpu = 500
        memory = 256

        network {
          mbits = 10

          port "grpc" {}
        }
      }
    }
  }
}
