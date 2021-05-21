package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	ip    string
	host  string
	ns    bool
	mx    bool
	txt   bool
	cname bool
)

func init() {
	flag.StringVar(&ip, "ip", "", "IP address for DNS operation")
	flag.StringVar(&host, "host", "", "Host address for DNS operations")
	flag.BoolVar(&ns, "ns", false, "Host name server lookup")
	flag.BoolVar(&mx, "mx", false, "Host domain server lookup")
	flag.BoolVar(&txt, "txt", false, "Host domain TXT lookup")
	flag.BoolVar(&cname, "cname", false, "Host CNAME lookup")
}

type lsdns struct {
	res *net.Resolver
	ctx context.Context
}

func newLsdns() *lsdns {
	return &lsdns{net.DefaultResolver, context.Background()}
}

func (ls *lsdns) reverseLkp(ip string) error {
	names, err := ls.res.LookupAddr(ls.ctx, ip)
	if err != nil {
		warn(err)
		return err
	}

	prln()
	fmt.Println("Reverse Lookup")
	prln()
	for _, name := range names {
		fmt.Println(name)
	}
	fmt.Println("")
	return nil
}

func (ls *lsdns) hostLkp(host string) error {
	addrs, err := ls.res.LookupHost(ls.ctx, host)
	if err != nil {
		warn(err)
		return err
	}
	prln()
	fmt.Println("Host Lookup")
	prln()
	for _, addr := range addrs {
		fmt.Printf("%-30s%-20s\n", host, addr)
	}
	fmt.Println("")
	return nil
}

func (ls *lsdns) nsLkp(host string) error {
	nss, err := ls.res.LookupNS(ls.ctx, host)
	if err != nil {
		warn(err)
		return err
	}
	prln()
	fmt.Println("NS Lookup")
	prln()
	for _, ns := range nss {
		fmt.Printf("%-25s%-20s\n", host, ns.Host)
	}
	fmt.Println("")
	return nil
}

func (ls *lsdns) mxLkp(host string) error {
	mxs, err := ls.res.LookupMX(ls.ctx, host)
	if err != nil {
		warn(err)
		return err
	}
	prln()
	fmt.Println("MX Lookup")
	prln()
	for _, mx := range mxs {
		fmt.Printf("%-17s%-11s\n", host, mx.Host)
	}
	fmt.Println("")
	return nil
}

func (ls *lsdns) txtLkp(host string) error {
	txts, err := ls.res.LookupTXT(ls.ctx, host)
	if err != nil {
		warn(err)
		return err
	}
	prln()
	fmt.Println("TXT Lookup")
	prln()
	for _, txt := range txts {
		fmt.Println(txt)
	}
	fmt.Println("")
	return nil
}

func (ls *lsdns) cnameLkp(host string) error {
	cname, err := ls.res.LookupCNAME(ls.ctx, host)
	if err != nil {
		warn(err)
		return err
	}
	prln()
	fmt.Println("CNAME Lookup")
	prln()
	fmt.Println(cname)
	fmt.Println("")
	return nil
}

func warn(err error) {
	fmt.Fprintf(os.Stderr, "warning: %v", err)
}

func prln() { fmt.Println("--------------------------------------------") }

func main() {
	flag.Parse()
	ls := newLsdns()

	switch {
	case ip != "":
		ls.reverseLkp(ip)
	case host != "":
		switch {
		case ns:
			ls.nsLkp(host)
		case mx:
			ls.mxLkp(host)
		case txt:
			ls.txtLkp(host)
		case cname:
			ls.cnameLkp(host)
		default:
			ls.hostLkp(host)
		}
	default:
		fmt.Println("flag ip or host must be provided")
	}
}
