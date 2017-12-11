package atmin

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"net/http"
)

type HTTPExecutor struct {
	Addr string
	TLS  bool
}

func (m Minimizer) ExecuteHTTP(addr string, useTLS bool) Minimizer {
	m.ex = &HTTPExecutor{Addr: addr, TLS: useTLS}
	m.out = m.ex.Execute(m.in)

	return m
}

func (ex *HTTPExecutor) Execute(in []byte) []byte {
	var conn net.Conn
	var err error

	r, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(in)))
	if err != nil {
		return []byte(err.Error())
	}

	if ex.TLS {
		conn, err = tls.Dial("tcp", ex.Addr, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		conn, err = net.Dial("tcp", ex.Addr)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer conn.Close()

	conn.Write(in)

	resp, err := http.ReadResponse(bufio.NewReader(conn), r)
	if err != nil {
		return []byte(err.Error())
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	return buf.Bytes()
}
