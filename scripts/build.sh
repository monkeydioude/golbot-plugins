#!/bin/bash

pkg=github.com\\/monkeydioude\\/golmods
imports=
commands=

for i in $(cd pkg && ls -d *);
do
    imports=$imports"\\\"$pkg\\/pkg\\/$i\\\"\n\t" 
    # commands=$commands$i.AddCommand\(g,\ g.NewCache\(\"$i\"\)\),"\n\t\t"
    commands=$commands$i.AddCommand\(g\),"\n\t\t"
done;

sed "s/#MODS#/$imports/;s/#ADD_COMMAND#/$commands/;" plugins.go.tpl > plugins.go
