#!/bin/bash
cd ../aozorabunko
for i in cards/*/files/; do mkdir -p ../azb_serverfiles/$i; cp --update $i*.{html,png} ../azb_serverfiles/$i; done && cd -
cp ../aozorabunko/index_pages/list_person_all_extended_utf8.zip index_pages/
date -u -R
