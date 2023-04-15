#!/usr/bin/env bash
#sudo apt install -y software-properties-common
#sudo add-apt-repository -y ppa:jonathonf/ffmpeg-3
#sudo apt update
#sudo apt install -y ffmpeg screen

wget http://deployer:eZKpb4uR9cEW9A7tvQF2ug3LCNZxxpzysuDYYBzbUn@teamcity.misc.vee.bz/repository/download/Migrator_Build/65/migrator/migrator -O migrator

chmod +x migrator

/usr/bin/screen -S migration -dm ./migrator -legacy-aphrodite-uri=http://api.aphrodite.sla.bz/ \
-legacy-aphrodite-token=Mf9MAsk7wFT4bco34ELepPvTa7NpzY6omDrCvwM \
-aphrodite-uri=http://api.aphrodite.vee.bz \
-aphrodite-token=helloWorld \
-minerva-dsn=minerva.vee.bz:8081 \
-minerva-read-token=read \
-minerva-write-token=write \
-dry-run=0 -max-transfers=0 -verbose=0 -skip-unassociated=1 -routines=8
