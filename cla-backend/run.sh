#!/bin/bash
#gunicorn cla.routes:__hug_wsgi__ -b 0.0.0.0:8080 --log-level debug
# Temporary
hug -f cla/routes.py -p 8080
