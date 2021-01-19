
name="${1}"

cd benchmarks
go test -bench "${name}" -benchmem -o pprof/test.bin  -cpuprofile pprof/cpu.out .
go tool pprof --svg pprof/test.bin pprof/cpu.out > pprof/test.svg
open pprof/test.svg
