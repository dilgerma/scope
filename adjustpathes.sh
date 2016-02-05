#!/bin/bash
find . -name '*.go'  | xargs sed -i "bak" 's%github.com/weaveworks/scope%github.com/dilgerma/scope%g'
find . -name *.gobak | xargs rm 
