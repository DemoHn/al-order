#!/bin/sh
run () {  
  docker-compose stop && docker-compose up  
}

run