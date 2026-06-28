#!/bin/bash

#Run this script **before** service workers are added to the dist folder!  
#There is a rudimentary guard against caching the service worker itself but
#it's best to generate the manifest before the service workers are even present.

cd dist

FILE="manifest.js"

echo "$(dir -1 --file-type)" > $FILE
echo "$(dir -1 stylesheets/*)" >> $FILE
echo "$(dir -1 scripts/*)" >> $FILE
echo "index_pages/list_person_all_extended_utf8.zip" >> $FILE

sed -i '/.*sw.js/d' $FILE
sed -i '/.*\/$/d' $FILE
sed -i -e 's/.*/"\/\0",/' $FILE

sed -i '1i export const\ FILES\ =\ [' $FILE

echo "]" >> $FILE
