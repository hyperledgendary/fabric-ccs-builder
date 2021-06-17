build:
    go build -o bin ./cmd/build/
    go build -o bin ./cmd/detect/
    go build -o bin ./cmd/release/

test: build
    rm -rf ./_test/release && rm -rf ./_test/output
    ./bin/detect ./_test/src/ ./_test/md/ && echo "== Detect OK =="
    ./bin/build ./_test/src/ ./_test/md/ ./_test/output/ && echo "== Build OK =="
    ./bin/release ./_test/output/ ./_test/release/ && echo "== Release OK =="
    cp -r ./bin ../../hyperledger/fabric-samples/external-chaincode/builder/