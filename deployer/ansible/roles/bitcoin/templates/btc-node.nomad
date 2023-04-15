job "btcd" {
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

  update {
    stagger = "5s"
    max_parallel = 1
  }

  group "node" {
    count = 1

    restart {
      attempts = 6
      interval = "1m"
      delay = "10s"
      mode = "delay"
    }

    task "node" {
      driver = "raw_exec"

      config {
        command = "bitcoin-0.14.2/bin/bitcoind"
        args = [
          "-disablewallet",
          "-conf=/etc/bitcoin/bitcoind.conf"
        ]
      }

      artifact {
        source = "https://bitcoin.org/bin/bitcoin-core-0.14.2/bitcoin-0.14.2-x86_64-linux-gnu.tar.gz"
      }

      service {
        name = "btc-rpc"
        port = "rpc"
        tags = [
          "btc",
          "rpc"
        ]
      }

      service {
        name = "btc-node"
        port = "node"
        tags = [
          "btc",
          "node"
        ]
      }

      resources {
        cpu = 1000
        memory = 8124

        network {
          mbits = 100

          port "rpc" {
            static = "8332"
          }

          port "node" {
            static = "8333"
          }
        }
      }
    }
  }
}

