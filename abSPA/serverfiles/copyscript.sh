#!/bin/bash
cd ../aozorabunko
git pull
for i in cards/*/files/; do mkdir -p ../azb_serverfiles/cards/002265/files/; cp --update cards/002265/files/*.{html,png} ../azb_serverfiles/cards/002265/files/; done && cd -
cp ../aozorabunko/index_pages/list_person_all_extended_utf8.zip index_pages/
date -u -R
