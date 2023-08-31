# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

FROM public.ecr.aws/lambda/python:3.7

# Copy requirements.txt
ADD requirements.txt requirements-test.txt ${LAMBDA_TASK_ROOT}

# Copy function code
ADD *.py .serverless-wsgi serverless-authorizer.yml ${LAMBDA_TASK_ROOT}
ADD /cla ${LAMBDA_TASK_ROOT}/cla

# Install the specified packages
RUN pip install -r requirements.txt

# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
CMD [ "wsgi_handler.handler" ]
