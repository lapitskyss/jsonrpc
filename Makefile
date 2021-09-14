test:
	go test  -v .

bench:
	go test -bench=.

bench_1:
	go test -bench=BenchmarkServeHTTP -benchmem -benchtime=1000x -cpu=1

prof:
	go test -bench=BenchmarkServeHTTP -memprofile=mem.out -cpuprofile=cpu.out -outputdir=optimization

http_prof:
	go tool pprof -http :8082 jsonrpc.test ./optimization/mem.out

alloc:
	go tool pprof -http :8082 --alloc_objects jsonrpc.test ./optimization/mem.out
