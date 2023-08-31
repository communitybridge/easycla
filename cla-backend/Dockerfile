# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

FROM public.ecr.aws/lambda/python:3.7

# Copy requirements.txt
COPY requirements.txt ${LAMBDA_TASK_ROOT}
# Copy function code
COPY *.py .serverless-wsgi serverless-authorizer.yml ${LAMBDA_TASK_ROOT}
# Add the CLA library folder
ADD /cla ${LAMBDA_TASK_ROOT}/cla
ADD /helpers ${LAMBDA_TASK_ROOT}/helpers
ADD /auth ${LAMBDA_TASK_ROOT}/auth

# Install the specified packages
RUN pip install --upgrade pip && pip install -r requirements.txt

# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
CMD [ "wsgi_handler.handler" ]
