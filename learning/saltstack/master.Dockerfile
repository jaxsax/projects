FROM ubuntu:20.04

RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
RUN apt-get update -y && apt-get install apt-utils wget gnupg -y
RUN wget -q https://repo.saltstack.com/py3/ubuntu/20.04/amd64/latest/SALTSTACK-GPG-KEY.pub && \
    apt-key add SALTSTACK-GPG-KEY.pub && \
    rm SALTSTACK-GPG-KEY.pub
RUN echo "deb http://repo.saltstack.com/py3/ubuntu/20.04/amd64/latest focal main" > /etc/apt/sources.list.d/saltstack.list
RUN apt-get update -y && \
    apt-get install salt-master salt-minion salt-ssh salt-syndic salt-cloud salt-api -y && \
    apt-get clean all

RUN sed -i "s|#auto_accept: False|auto_accept: True|g" /etc/salt/master

ENTRYPOINT ["salt-master", "-l", "debug"]




