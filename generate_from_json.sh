#!/bin/bash

gojson -input .generate/bulk_input_header.json -name BulkHeader -pkg elk
gojson -input .generate/bulk_input_data.json -name BulkMessage -pkg elk
