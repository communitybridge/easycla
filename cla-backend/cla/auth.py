"""
auth.py contains all necessary objects and functions to perform authentication and authorization.
"""
import os
import requests
from jose import jwt
import cla

auth0_base_url = os.environ.get('AUTH0_DOMAIN', '')
auth0_username_claim = 'https://sso.linuxfoundation.org/claims/username'
auth0_audience_claim = 'hquZHO8JNsaIScoayPtCS5VELdn7TnVq'
algorithms = ['RS256']

class AuthError(Exception):
    def __init__(self, response):
        self.response = response

def get_auth_token(headers):
    """
    Obtains the Access Token from the Authorization Header
    """
    auth = headers.get('Authorization')
    if not auth:
        auth = headers.get('AUTHORIZATION')
    if not auth:
        raise AuthError('missing authorization header')

    parts = auth.split()

    if parts[0].lower() != 'bearer':
        raise AuthError({'authorization header must begin with \"Bearer\"'})
    elif len(parts) == 1:
        raise AuthError('token not found')
    elif len(parts) > 2:
        raise AuthError('authorization header must be of the form \"Bearer token\"')

    return parts[1]

def authenticate_user(request):
    """
    Determines if the Access Token is valid
    """
    token = get_auth_token(request.headers)
    try:
        jwks_url = os.path.join('https://', auth0_base_url, '.well-known/jwks.json')
        jwks = requests.get(jwks_url).json()
    except requests.exceptions.RequestException as e:
        raise AuthError('unable to fetch well known jwks')

    try:
        unverified_header = jwt.get_unverified_header(token)
    except jwt.JWTError as e:
        cla.log.error(e)
        raise AuthError('unable to decode claims')

    rsa_key = {}
    for key in jwks["keys"]:
        if key["kid"] == unverified_header["kid"]:
            rsa_key = {
                "kty": key["kty"],
                "kid": key["kid"],
                "use": key["use"],
                "n": key["n"],
                "e": key["e"]
            }
    if rsa_key:
        try:
            payload = jwt.decode(
                token,
                rsa_key,
                algorithms=algorithms,
                audience=auth0_audience_claim,
                # issuer="https://"+auth0_base_url+"/"
                options={
                    'verify_at_hash': False
                }
            )
        except jwt.ExpiredSignatureError as e:
            cla.log.error(e)
            raise AuthError('token is expired')
        except jwt.JWTClaimsError as e:
            cla.log.error(e)
            raise AuthError('incorrect claims')
        except Exception as e:
            cla.log.error(e)
            raise AuthError('unable to parse authentication')

        username = payload.get(auth0_username_claim)
        if username is None:
            raise AuthError('username not found')

        payload['username'] = username

        return payload

    raise AuthError({"code": "invalid_header",
                    "description": "Unable to find appropriate key"})
