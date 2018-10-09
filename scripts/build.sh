#!/bin/bash

pkg=github.com\\/monkeydioude\\/golmods
imports=

for i in $(cd pkg && ls -d *);
do
    imports=$imports"_ \\\"$pkg\\/pkg\\/$i\\\"\n\t" 
done;

sed "s/#MODS#/$imports/;" plugins.go.tpl > plugins.go
