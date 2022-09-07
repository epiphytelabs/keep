package server

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

const (
	Root = "keep."
)

func resolve(w dns.ResponseWriter, r *dns.Msg) {
	fmt.Printf("resolve type=%s name=%s\n", dns.TypeToString[r.Question[0].Qtype], r.Question[0].Name)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Question[0].Qtype {
	case dns.TypeA:
		m.Answer = append(m.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
			A:   net.IPv4(127, 0, 0, 1),
		})
	}

	if r.IsTsig() != nil {
		if w.TsigStatus() == nil {
			m.SetTsig(r.Extra[len(r.Extra)-1].(*dns.TSIG).Hdr.Name, dns.HmacMD5, 300, time.Now().Unix())
		} else {
			println("tsig status", w.TsigStatus().Error())
		}
	}

	if err := w.WriteMsg(m); err != nil {
		println("write error", err)
	}
}
