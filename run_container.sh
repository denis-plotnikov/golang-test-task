#!/bin/bash

# TSERVER_HOST and TSERVER_HOST enviroment variables to set corresponding service values

REDIRECT_PORT=13000
SERVICE_PORT=50000
docker run --name golang-test-task -e TSERVER_HOST="0.0.0.0" -e TSERVER_PORT=${SERVICE_PORT} -p ${REDIRECT_PORT}:${SERVICE_PORT} --rm golang-test-task
