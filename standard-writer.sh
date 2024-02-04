#!/bin/bash

if [ "$1" = "-date" ] && [ -n "$2" ]; then
    if ! date -d "$2" &> /dev/null; then
        echo "Invalid date format: $2"
        exit 1
    fi
    target_date=$(date -d "$2" '+%Y-%m-%d')
else
    target_date=$(date '+%Y-%m-%d')
fi

for i in {0..6}; do
    day_to_process=$(date -d "$target_date - $i days" '+%Y-%m-%d')
    ./writer -date "$day_to_process"
    sleep 1m
done
