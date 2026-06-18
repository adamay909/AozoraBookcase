#!/bin/bash

aws cloudfront create-invalidation --distribution-id=EVCQBBG3MBV21 --paths "/aozorabookcase.html" "/bookcase.wasm" "/wasm_exec.js" "/ebooks.css" "/readingpane.css"

