#!/bin/sh

# Define a timestamp function
timestamp() {
  date +"updated -> %D %T" # current time
}

R=`timestamp`

DATA=`find .. -name '*.go' | xargs wc -l`

echo "<|--------------|>\n\t$R\n$DATA\n<|--------------|>\n">> codestat.log

# date +"updated -> %D %T" >> codestat.log
# find .. -name '*.go' | xargs wc -l >> codestat.log