package atmin

import (
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"time"
)

type NetExecutor struct {
	Addr string
	TLS  bool
}

// TODO(amlw): there must be a better way...
func (m Minimizer) ExecuteNet(addr string, useTLS bool) Minimizer {
	m.ex = &NetExecutor{Addr: addr, TLS: useTLS}
	m.out = m.ex.Execute(m.in)

	return m
}

func (ex *NetExecutor) Execute(in []byte) []byte {
	var conn net.Conn
	var err error

	// bail out early if we know the request is invalid
	if !bytes.Contains(in, []byte("\r\n\r\n")) && !bytes.Contains(in, []byte("\n\n")) {
		return []byte("invalid HTTP request")
	}

	if ex.TLS {
		conn, err = tls.Dial("tcp", ex.Addr, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		conn, err = net.Dial("tcp", ex.Addr)
		if err != nil {
			// TODO(amlw): should executors return errors?
			log.Fatal(err)
		}
	}
	defer conn.Close()

	// set deadlines just in case
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	conn.Write(in)
	var b = make([]byte, 4096)
	// TODO(amlw): why the hell must we Read() instead of ioutil.ReadAll()
	conn.Read(b)
	return b
}
