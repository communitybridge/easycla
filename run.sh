#!/bin/bash
gunicorn cla.routes:__hug_wsgi__ -b 0.0.0.0:8080
