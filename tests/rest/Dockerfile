FROM python:3.7-alpine

WORKDIR /app
COPY . .
RUN pip install -r requirements.freeze.txt


CMD pytest --alluredir=allure_result_folder -v
