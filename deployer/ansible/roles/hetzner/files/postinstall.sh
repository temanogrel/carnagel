#!/usr/bin/env bash

cat << EOF > /etc/apt/sources.list
deb mirror://mirrors.ubuntu.com/mirrors.txt xenial main restricted universe multiverse
deb mirror://mirrors.ubuntu.com/mirrors.txt xenial-backports main restricted universe multiverse
deb mirror://mirrors.ubuntu.com/mirrors.txt xenial-updates main restricted universe multiverse
deb mirror://mirrors.ubuntu.com/mirrors.txt xenial-security main restricted universe multiverse
EOF

apt-get update
apt-get install -y python2.7
mkdir -p /root.ssh
echo "rsa AAAAB3NzaC1yc2EAAAABJQAAAgEAmBADizn+tSJabyD2w4vFrcVGvFuAh0u7Rdvg9+Uwzgojdeq8JHr3b/D9hat5VFGXDsG+/Jr0o5xZpuBKXmluU6OOcnddXLAbKd/pZ3ZiXmycYsRDyxgEAlTRnVMn4vVdDC6+Xy6daWkznShrP2K16BIH87SFKgA3T2PVic6TvCFFRBIrS+1zSRAkAKNnRzOruVWiB2b5yYnJ5lSI9NbP+KCapBJVxg/KzlpnoSC5/ELGXQftM9mc8UcKcAfpvGzhrcN+s3mVzqKL0rG4byHT+S1SeOR8NKBzIPefn51cTYinvirvsNd7KuSGpXKM3v7wXnHPU0l1CEqgiSstmr6X82+QehHnTf3tv3Lxjp1nVaBgdlu+w3JqzTkrFyeEAglWY7utLf+wiQCqJ0X2jrkXEd323zFloAyFFYfu9Z0kGfFt7DWPqUMRod/3XL69nQGeqnIp0kbrdBODLuw1xXJXIy4n+sUKc+XCl0D6ojMyJoaB2BrXw55NCRSakVTgML+B7LFDlnd/UBbYXMbjkKkYb6B+6TgqD88A0WPHuJc/tO9pSQL8nwSGNckG/cQ/HgyJFquPs1UxINmIVZtUocCyRsYFMw3AzNdU5ejXZEIIzZLIyAEfN+jRO/uYcctr//tTaxOaPsZWZbk9xmU9TkfZSmjB5CJqoC4Gv92toD042m0=" > /root/.ssh/authorized_keys
