# FROM ubuntu:20.04
FROM alpine:3.15.5

ENV TF_VERSION 0.14.3

ENV KUBESPRAY_REPO_URL https://github.com/sjoh0704/kubespray.git

ENV BRANCH operator

ENV DEBIAN_FRONTEND noninteractive

# 필요 패키지 설치 
RUN apk --no-cache --update-cache add \
    #   libressl openssh-client sshpass python3 py3-pip
    libressl openssh-client sshpass python3 ca-certificates

# terraform binary install
RUN apk --no-cache --update-cache --virtual .build-deps add \ 
    curl unzip \
    && curl -O https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip \
    && unzip terraform_${TF_VERSION}_linux_amd64.zip \
    && mv terraform /usr/local/bin/ \
    && rm -rf terraform_${TF_VERSION}_linux_amd64.zip \
    && terraform version \
    && apk del .build-deps  

# kubespray git repo clone 
RUN apk --no-cache --update-cache --virtual .build-deps add git \
    && git clone -b ${BRANCH} --single-branch --depth 1 ${KUBESPRAY_REPO_URL} \
    && rm -rf /kubespray/contrib/terraform/aws/terraform.tfvars \
    && apk del .build-deps

WORKDIR kubespray

# script file 추가 
COPY setting.sh .
COPY execute.sh .
COPY provision.sh .
COPY destroy.sh .
COPY install.sh .

# kubespray requirements 설치
RUN apk --no-cache --update-cache --virtual .build-deps add \
    gcc make python3-dev py3-pip linux-headers musl-dev libffi-dev openssl-dev \
    && pip3 install -r requirements.txt \
    && apk del .build-deps

# terraform init 포함
RUN terraform -chdir=/kubespray/contrib/terraform/aws/ init
# RUN apk add upx \
#     && upx /usr/local/bin/terraform \
#     && upx /usr/bin/

# 특이사항
