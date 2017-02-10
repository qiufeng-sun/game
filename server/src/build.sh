#!/bin/bash

#find dest dir
p=$(dirname $0)
echo "dir:$p"

#go to dest dir && build
out="w.run"
(cd $p && cd .. && . gvp && cd - && pwd && go build -o $out world && echo "build ok! out: $out")
