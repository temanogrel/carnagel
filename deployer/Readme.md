# Provision new Hetzner machines
The process to setup and integrate new machines from Hetzner is divided into three setp
The first two steps require a different inventory then step three (and all other tasks).

New machines should be but into a temporary group until all steps were run successfully. Afterward they shuold be moved to their corresponding group.

## Install OS from Rescue system
- `ansible_user=root` 
- `ansible_password=` Check Google spreadsheat or Hetzner UI
- `ansible_python_interpreter=/usr/bin/python2.7`
- `ansible-playbook --extra-vars="step=install_os" ansible/hetzner.yml`

This step also install python 2.7 which will then be located at `/usr/bin/python2.7`

## Base setup for freshly installed OS
- `ansible_user=root` 
- `ansible_password=` Check Google spreadsheat or Hetzner UI
- `ansible_python_interpreter=/usr/bin/python2.7`
- `ansible-playbook --extra-vars="step=init" ansible/hetzner.yml`

This step creates the default user `ansible`, adds it to the `sudo` group and restores the original `/etc/apt/sources.list` from Hetzner.

## Harden sshd and setup internal network
- `ansible_python_interpreter=/usr/bin/python2.7`
- `internal_ip=10.0.x.x` Check Google Spreadhsheet
- `internal_subnet=x.x.x.x` Check Google Spreadhsheet
- `internal_vlan=10.0.x.0` Check Google Spreadhsheet
- `ansible-playbook --private-key=XXX --extra-vars="step=network ansible_sudo_pass=bellanaleck" ansible/hetzner.yml`

We no longer use password for authentication but public key.
