// demo issues an HTTPS GET with a TC eBPF SYN-spoofing profile applied.
// Run on Linux with CAP_NET_ADMIN + CAP_BPF (or as root).
//
//	go build -o /tmp/demo ./tcp/cmd/demo
//	sudo /tmp/demo -url https://robinsamuel.dev/
//
// By default the BPF program attaches to the auto-detected default-route
// interface (overridable with -iface or the TCP_IFACE env var).
//
// Watch the SYN go out in another terminal:
//
//	sudo tcpdump -nn -vvv -i eth0 \
//	  'tcp[tcpflags] & tcp-syn != 0 and not tcp[tcpflags] & tcp-ack != 0'
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/nukilabs/tlsclient"
	"github.com/nukilabs/tlsclient/profiles"
	"github.com/nukilabs/tlsclient/tcp"
)

func main() {
	iface := flag.String("iface", "", "egress interface (default: auto-detect from /proc/net/route)")
	target := flag.String("url", "https://robinsamuel.dev/", "URL to fetch")
	flag.Parse()

	if *iface != "" {
		tcp.SetInterface(*iface)
	}

	client := tlsclient.New(profiles.Chrome(120),
		tlsclient.WithTCP(tcp.Windows),
		tlsclient.WithTimeout(15*time.Second),
	)

	fmt.Printf("GET %s with tcp.Windows ...\n", *target)
	resp, err := client.Get(*target)
	if err != nil {
		log.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read body: %v", err)
	}
	fmt.Printf("status: %s\n", resp.Status)
	fmt.Printf("body (%d bytes):\n%s\n", len(body), string(body))
}
