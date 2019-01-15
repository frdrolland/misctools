#!/bin/bash
kill $(ps -ef | grep misctools | grep -v grep | awk '{ print $2}')
