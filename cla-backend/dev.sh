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
    pip3 install -r requirements.txt &&\
    cat cla/config.py > cla_config.py &&\
  echo '======> setting up aws profile..'
    cd .. &&\
    serverless config credentials --provider aws --profile lf-cla --key ' ' --secret ' ' -s devS &&\
    cd cla-backend
  echo '======> installing dynamodb local..'
    sls dynamodb install -s 'local' &&\
  echo '======> installation has done. please run npm run add:user github|#######'

elif [ $1 = 'add:user' ]; then
  if [ "x" != "x$2" ]; then
  echo '======> creating permission in local db'
   aws dynamodb put-item \
    --table-name "cla-local-user-permissions" \
    --item '{ "user_id": { "S": "'$2'" }, "projects": { "SS": ["a09J000000KHoZVIA1","a09J000000KHoayIAD"] } }' \
    --profile lf-cla --region "us-east-1" \
    --endpoint-url http://localhost:8000
  echo '======> done!'
  else
    echo '======> user id is required!'
  fi

elif [ $1 = 'start:lambda' ]; then
  echo '======> activating virtual enviroment'
    . ~/.env/lf-cla/bin/activate &&\
  echo '======> running local lambda server'
    sls wsgi serve -s 'local'

elif [ $1 = 'start:dynamodb' ]; then
  echo '======> running local dynamodb server'
    sls dynamodb start -s 'local'

elif [ $1 = 'start:s3' ]; then
  echo '======> running local s3 server'
    sls s3 start -s 'local'
  
else
  echo "option not valid"
  exit 0
fi