## Fuzzer Tools for goslp

Fuzzing helps find bugs which may have been missed by unit tests.  A differential fuzzer has been implemented for comparing the goslp v1parser and slp-validate javascript library.

### Fuzzing v1parser.ParseSLP()

The v1parser has been setup in `../v1parser/v1_parer_fuzz.go` in order to run with dvyukov/go-fuzz.  A corpus file `corpus.tar.gz` has been included which was generated from previous fuzzing campaigns.  Extract the contents of this file to `v1parser/corpus`.


The following commands will run the fuzzer.

```
# start http server for getting slp-validate result
npm i
node ./server.js

# extract initial seed corpus to v1parser directory
$ cd v1parser
$ tar -xzvf ../fuzz/corpus.tar.gz

# build and start fuzzer
$ go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build
$ go-fuzz-build
$ go-fuzz --procs 4
```
