#!/bin/sh
docker build -t yadro_test .
docker run -it yadro_test test_file_1.txt