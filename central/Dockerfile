FROM ubuntu:latest

# Install OS dependencies
RUN apt-get update \
    && apt-get install -y git locales python3-dev python3-pip libxml2-dev libxslt1-dev zlib1g-dev \
    && apt-get clean

# Setup ssh
RUN mkdir /root/.ssh;
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN ssh-keyscan git.misc.vee.bz >> /root/.ssh/known_hosts

COPY docker/ssh_config /root/.ssh/config
COPY docker/deploy_rsa /root/.ssh/deploy_rsa
COPY docker/gunicorn_rsa /root/.ssh/gunicorn_rsa

RUN chmod 600 /root/.ssh/deploy_rsa
RUN chmod 600 /root/.ssh/gunicorn_rsa

# Setup the locale
RUN locale-gen en_US.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8

# Add the code to the application
ADD ./central /opt/application/central
ADD ./requirements.txt /opt/application/requirements.txt
ADD ./application /opt/application/application

# Change work dir
WORKDIR /opt/application

# Install application dependencies
RUN eval $(ssh-agent -s) \
    && ssh-add /root/.ssh/deploy_rsa \
    && /usr/bin/pip3 install -r requirements.txt --upgrade
