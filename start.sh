#!/bin/sh
run () {
  # run mysql & redis (to host network)
  docker-compose stop && docker-compose up
  # create container  
}

run