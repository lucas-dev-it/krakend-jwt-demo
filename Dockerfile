FROM python:3.9.0b5-alpine3.12 as builder
WORKDIR /tmp/builder
COPY ./krakend .
