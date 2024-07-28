#!/bin/bash
files=`ls schema/*.graphql`
for file in $files
do
    echo "Formatting $file"
    format-graphql --write true --sort-arguments true --sort-definitions true --sort-fields false $file
done
