#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings -dbname=go -dbuser=tsawlergo -dbpass=1234 -cache=false -production=false


 