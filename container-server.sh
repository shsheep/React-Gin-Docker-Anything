#!/bin/bash

while true; do
   inotifywait -e create --format '%w%f' "./input" | while read FILE; do
      go run $FILE -json_output_file output/json_output.txt
   done
done
