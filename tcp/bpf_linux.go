package tcp

// bpf_linux.go drives bpf2go to compile bpf/syn_rewrite.c into Go-embeddable
// objects. Running `go generate ./tcp/...` produces:
//
//   synrewrite_bpfel.go  + synrewrite_bpfel.o   (little-endian: amd64, arm64, ...)
//   synrewrite_bpfeb.go  + synrewrite_bpfeb.o   (big-endian:    s390x, ...)
//
// Prerequisites on the build host:
//   - clang and llvm-strip in $PATH
//   - linux-libc-dev / kernel UAPI headers under /usr/include/linux
//
// The compiled .o files contain a relocatable BPF object the loader (step 3)
// reads via the generated loadSynrewrite() helper.

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -Wall -Werror -nostdinc" -target amd64,arm64 -output-dir . synrewrite ./bpf/syn_rewrite.c -- -I./bpf/headers
