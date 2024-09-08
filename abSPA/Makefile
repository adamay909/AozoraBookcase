build:
	GOOS=js GOARCH=wasm go build -o serverfiles/bookcase.wasm

deploy: build
	cp -f /usr/lib/go/misc/wasm/wasm_exec.js serverfiles/ &&\
	cp -f serverfiles/* ${AZSPADIR}/

