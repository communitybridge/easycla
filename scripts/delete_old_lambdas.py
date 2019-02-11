import boto3

# python3 delete_old_lambdas.py
def clean_old_lambda_versions():
    session = boto3.Session(profile_name='lf-cla')
    client = session.client('lambda', region_name='us-east-1')
    functions = client.list_functions()['Functions']
    for function in functions:
        versions = client.list_versions_by_function(FunctionName=function['FunctionArn'])['Versions']
        for version in versions:
            arn = version['FunctionArn']
            if version['Version'] != function['Version'] and arn.find('cla-') != -1:
                print('delete_function(FunctionName={})'.format(arn))
                # client.delete_function(FunctionName=arn)  # uncomment me once you've checked
                


if __name__ == '__main__':
    clean_old_lambda_versions()