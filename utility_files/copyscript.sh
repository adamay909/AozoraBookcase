#!/bin/bash

cd ../aozorabunko

git pull

for i in cards/*/files; do mkdir -p ../azb_serverfiles/${i}/; cp --update ${i}/*.{html,png} ../azb_serverfiles/${i}/; done 

cd -

cp ../aozorabunko/index_pages/list_person_all_extended_utf8.zip index_pages/

aws s3 sync --exclude="*.sh" . s3://aozorabookcase

aws cloudfront create-invalidation --distribution-id=EVCQBBG3MBV21 --paths "/index_pages/*"

date -u -R
