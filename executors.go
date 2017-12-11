package atmin

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type NetExecutor struct {
	Addr string
	TLS  bool
}

func (m Minimizer) ExecuteNet(addr string, useTLS bool) Minimizer {
	m.ex = &NetExecutor{Addr: addr, TLS: useTLS}
	m.out = m.ex.Execute(m.in)

	return m
}

func (ex *NetExecutor) Execute(in []byte) []byte {
	var conn net.Conn
	var err error

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

	// set deadlines just in case
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	conn.Write(in)

	var b = make([]byte, 4096)
	conn.Read(b)

	return b
}

type HTTPExecutor struct {
	URL *url.URL
}

func (m Minimizer) ExecuteHTTP(url *url.URL) Minimizer {
	m.ex = &HTTPExecutor{URL: url}
	m.out = m.ex.Execute(m.in)

	return m
}

func (ex *HTTPExecutor) Execute(in []byte) []byte {
	r, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(in)))
	if err != nil {
		return []byte(err.Error())
	}
	r.URL = ex.URL
	r.RequestURI = ""

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return []byte(err.Error())
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	return buf.Bytes()
}
