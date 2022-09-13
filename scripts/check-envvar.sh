#!/bin/bash

param=$1

if [ -z "${!param}" ] ; then
  echo "Error: Environment variable $1 is undefined."
  exit 1
fi
