build:
	GOOS=js GOARCH=wasm go build -o serverfiles/bookcase.wasm
<<<<<<< HEAD

deploy: build
	if [[ "${AZSPADIR}x" == "x" ]]; 	then 	echo "Specify spa's root folder using the environment variable AZASPADIR"; exit 1; fi
	cp -f /usr/lib/go/misc/wasm/wasm_exec.js serverfiles/ &&\
=======
	cp -f /usr/lib/go/misc/wasm/wasm_exec.js serverfiles/ 

deploy: build
	if [[ "${AZSPADIR}x" == "x" ]]; then echo "Specify root folder of SPA using the environment variable AZASPADIR"; exit 1; fi
>>>>>>> 1c37364 (edit Makefile)
	cp -f serverfiles/* ${AZSPADIR}/

