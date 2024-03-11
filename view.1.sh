#!/usr/bin/env bash

if [ "$#" -lt 2 ] || [ "$#" -gt 3 ]; then
    echo "Usage: $0 <option1> <option2> [option3]"
    exit 1
fi

if [ -e annotated/$1*$2*.png ]; then
    open annotated/$1*$2*.png
fi

if [ "$#" -eq 3 ] && [ -e annotated/$1*$3*.png ]; then
    open annotated/$1*$3*.png
fi
