FROM python:3.7.4

WORKDIR /usr/src/app
RUN mkdir cla logs

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONBUFFERED 1

#install requirements
COPY ./requirements.txt .
RUN pip install -r requirements.txt

COPY ./entrypoint.sh /usr/src/app/entrypoint.sh
RUN chmod +x /usr/src/app/entrypoint.sh

COPY . /usr/src/app
