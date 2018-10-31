#!/bin/bash

if [ $1 = 'install' ]; then
  echo '======> installing npm dependencies..'
    npm install &&\
    npm i -g serverless &&\
  echo '======> installing python virtualenv..'
    pip3 install virtualenv &&\
  echo '======> creating virtual enviroment..'
    virtualenv ~/.env/lf-cla &&\
  echo '======> activating virtual enviroment'
    . ~/.env/lf-cla/bin/activate &&\
  echo '======> installing python dependencies..'
    pip install -r requirements.txt &&\
    cat cla/config.py > cla_config.py &&\
  echo '======> setting up aws profile..'
    cd .. &&\
    serverless config credentials --provider aws --profile lf-cla --key ' ' --secret ' ' -s devS &&\
    cd cla-backend
  echo '======> installing dynamodb local..'
    sls dynamodb install -s dev &&\
  echo '======> installation has done.'

elif [ $1 = 'start:lambda' ]; then
  echo '======> activating virtual enviroment'
    . ~/.env/lf-cla/bin/activate &&\
  echo '======> running local lambda server'
    AWS_PROFILE=lf-cla sls wsgi serve -s dev

elif [ $1 = 'start:dynamodb' ]; then
  echo '======> running local dynamodb server'
    AWS_PROFILE=lf-cla sls dynamodb start -s dev

elif [ $1 = 'start:s3' ]; then
  echo '======> running local s3 server'
    AWS_PROFILE=lf-cla sls s3 start -s dev
  
else
  echo "option not valid"
  exit 0
fi