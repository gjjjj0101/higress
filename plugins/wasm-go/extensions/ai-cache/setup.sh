version="1.0.23"

tinygo build -o main.wasm -scheduler=none -target=wasi -gc=custom -tags="custommalloc nottinygc_finalizer proxy_wasm_version_0_2_100" ./
docker build -t funbugjian/tianchi:$version -f Dockerfile .
docker push funbugjian/tianchi:$version