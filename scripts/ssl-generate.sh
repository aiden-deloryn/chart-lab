#!/bin/bash

openssl req -new -newkey rsa:2048 -nodes -keyout ssl/chartlab.key -out ssl/chartlab.csr
openssl x509 -req -days 365 -in ssl/chartlab.csr -signkey ssl/chartlab.key -out ssl/chartlab.crt