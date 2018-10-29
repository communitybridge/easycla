#!/bin/bash

if [ $1 = 'install' ]; then
  echo '======> installing npm dependencies..'
  npm install &&\
  echo '======> installing python virtualenv..'
  pip3 install virtualenv &&\
  echo '======> creating virtual enviroment..'
  virtualenv ~/.env/lf-cla &&\
  echo '======> activating virtual enviroment'
  . ~/.env/lf-cla/bin/activate &&\
  echo '======> installing python dependencies..'
  pip install -r requirements.txt
  cat cla/config.py > cla_config.py
  echo '======> installation has done, now run `npm run start:dev` to run local dev enviroment'

elif [ $1 = 'start:lambda' ]; then
  echo '======> activating virtual enviroment'
  . ~/.env/lf-cla/bin/activate &&\
  echo '======> running local lambda server'
  sls wsgi serve -s dev

elif [ $1 = 'start:dynamodb' ]; then
  echo '======> running local dynamodb server'
  sls dynamodb start -s dev

elif [ $1 = 'start:s3' ]; then
  echo '======> running local s3 server'
  sls s3 start -s dev
  
else
  echo "option not valid"
  exit 0
fi