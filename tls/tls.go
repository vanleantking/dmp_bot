package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
)

func main() {
	rootCertFile, er := ioutil.ReadFile("key.pem")
	if er != nil {
		panic(er.Error())
	}
	// rootCert, _ := rootCertFile.Read()
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootCertFile)
	if !ok {
		log.Fatal("failed to parse root certificate")
	}
	config := &tls.Config{RootCAs: roots}

	conn, err := tls.Dial("tcp", "localhost:8089", config)
	if err != nil {
		log.Fatal("log, ", err)
	}

	io.WriteString(conn, "Hello simple secure Server")
	conn.Close()
}
