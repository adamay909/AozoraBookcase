build:
	GOOS=js GOARCH=wasm go build -o serverfiles/bookcase.wasm
	cp -f /usr/lib/go/misc/wasm/wasm_exec.js serverfiles/ 

deploy: build
	if [[ "${AZSPADIR}x" == "x" ]]; 	then 	echo "Specify spa's root folder using the environment variable AZSPADIR"; exit 1; fi
	cp -f serverfiles/* ${AZSPADIR}/
	echo "done!"
