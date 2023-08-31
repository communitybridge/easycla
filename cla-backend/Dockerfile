# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

FROM public.ecr.aws/lambda/python:3.7

# Copy requirements.txt
COPY requirements.txt ${LAMBDA_TASK_ROOT}
# Copy function code
COPY *.py ${LAMBDA_TASK_ROOT}
COPY .serverless-wsgi ${LAMBDA_TASK_ROOT}
COPY serverless-authorizer.yml ${LAMBDA_TASK_ROOT}
# Add the CLA library folder
COPY /cla/ ${LAMBDA_TASK_ROOT}/cla/
COPY /helpers/ ${LAMBDA_TASK_ROOT}/helpers/
COPY /auth/ ${LAMBDA_TASK_ROOT}/auth/

# Install the specified packages
# RUN pip install --upgrade pip && pip install -r requirements.txt
RUN pip install -r requirements.txt

# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
CMD [ "wsgi_handler.handler" ]
