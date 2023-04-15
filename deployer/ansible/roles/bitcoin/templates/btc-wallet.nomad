job "btc-wallet" {
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

  group "wallet" {
    task "wallet" {
      driver = "docker"

      config {
        image = "halleburtin/btcwallet:{{ build_number }}"

        auth {
          email = "{{ lookup('env', 'DOCKER_EMAIL') }}"
          password = "{{ lookup('env', 'DOCKER_PASSWORD') }}"
          username = "{{ lookup('env', 'DOCKER_USERNAME') }}"
        }

        interactive = true
        network_mode = "host"
      }

      env {
        BTCD_HOST = "{{ lookup('env', 'BTCD_HOST') }}"
        BTCD_USER = "{{ lookup('env', 'BTCD_USER') }}"
        BTCD_PASS = "{{ lookup('env', 'BTCD_PASS') }}"
        WALLET_SEED = "{{ lookup('env', 'BTC_WALLET_SEED') }}"
        WALLET_RPC_CRT = "{{ lookup('env', 'BTC_WALLET_RPC_CRT') }}"
        WALLET_RPC_KEY = "{{ lookup('env', 'BTC_WALLET_RPC_KEY') }}"
      }

      service {
        name = "btc-wallet"
        port = "wallet"
        tags = [
          "btc",
          "wallet"
        ]
      }

      resources {
        cpu = 500
        memory = 2048

        network {
          mbits = 50

          port "wallet" {
            static = "8337"
          }
        }
      }
    }
  }
}