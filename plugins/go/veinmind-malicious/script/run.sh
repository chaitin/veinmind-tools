#!/bin/bash

source .env
docker-compose pull
docker-compose up -d --build --force-recreate