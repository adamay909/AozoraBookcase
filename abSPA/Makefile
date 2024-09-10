build:
	GOOS=js GOARCH=wasm go build -o serverfiles/bookcase.wasm

deploy: build
	if [[ "${AZSPADIR}x" == "x" ]]; 	then 	echo "Specify spa's root folder using the environment variable AZASPADIR"; exit 1; fi
	cp -f /usr/lib/go/misc/wasm/wasm_exec.js serverfiles/ &&\
	cp -f serverfiles/* ${AZSPADIR}/

