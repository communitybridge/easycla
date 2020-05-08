# EasyCLA REST API tests using tavern

This directory contains REST API tests built using
[Travern](https://github.com/taverntesting/tavern). The RestAPI test framework
based on [py.test](http://pytest.org/en/latest/). All the API tests are defined
in using the following file format: `test_*.tavern.yaml`. Reporting and
Visualization of tests are accomplished using
[allure](https://github.com/allure-framework/allure2) framework.

## Install dependencies to run tests in local

- Install python dependencies using pip

```bash
pip install -r requirements.freeze.txt
```

Optional step: Install [allure](https://github.com/allure-framework/allure2)
for reporting and visualization.

```bash
brew install allure

# OR

sudo apt-get update && \
sudo apt-get install software-properties-common -y  && \
sudo apt-add-repository ppa:qameta/allure -y && \
sudo apt-get update && \
sudo apt-get install allure -y
```

## Run Tests

```bash
export AUTH0_USERNAME=<username>
export AUTH0_PASSWORD=<password>
export AUTH0_CLIENT_ID=<client_id_of_eeasycla>
export AUTH0_AUDIENCE=<api_gateway_url>
export API_URL=<cla_api_url>
export STAGE=<dev_or_staging>

# Run tests
tavern-ci test_project_management_console.tavern.yaml --alluredir=allure_result_folder -v --tb=short

# OR

pytest --alluredir=allure_result_folder -v 

# Visualize tests with allure
allure serve allure_result_folder
```

## Run Tests in Docker

- Build docker image

```bash
docker build -t tavern-tests .
```

- Run tests in container

```bash
# Auth0 ID token
export AUTH_TOKEN="<TOKEN>"

export API_URL="https://api.staging.lfcla.com"
export STAGE="staging"

docker run --rm -it -e AUTH_TOKEN -e API_URL tavern-tests 
```

## List tests

```bash
pytest --collect-only
```

## Tavern yaml types and equivalent Python types

```python
_types = {
    "str": str,
    "int": int,
    "float": float,
    "number": None,
    "bool": bool,
    "map": dict,
    "seq": list,
    "timestamp": datetime.datetime,
    "date": datetime.date,
    "symbol": str,
    "scalar": None,
    "text": text,
    "any": object,
    "enum": str,
    "none": None
}
```
