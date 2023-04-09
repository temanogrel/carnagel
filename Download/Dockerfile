FROM ubuntu:latest

# Install OS dependencies
RUN apt-get update \
    && apt-get install -y git python3-dev python3-pip librtmp-dev libffi-dev nodejs libxml2-dev libxslt1-dev locales \
    && apt-get clean

# Setup ssh
RUN mkdir /root/.ssh;
RUN ssh-keyscan git.misc.vee.bz >> /root/.ssh/known_hosts

COPY docker/ssh_config /root/.ssh/config
COPY docker/deploy_rsa /root/.ssh/deploy_rsa
COPY docker/gunicorn_rsa /root/.ssh/gunicorn_rsa

RUN chmod 600 /root/.ssh/config
RUN chmod 600 /root/.ssh/deploy_rsa
RUN chmod 600 /root/.ssh/gunicorn_rsa

# Setup the locale
RUN locale-gen en_US.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8

RUN mkdir -p /opt/downloads

# Add the code to the application
ADD . /opt/application

# Change work dir
WORKDIR /opt/application

# Bug, cffi must be installed prior to installing requirements.
RUN /usr/bin/pip3 install cffi

# Install application dependencies
RUN /usr/bin/pip3 install -r requirements.txt --upgrade
