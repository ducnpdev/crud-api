package main

import (
	// "fmt"
	// "io"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("connect ok")
	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func main() {

	certs := x509.NewCertPool()

	pemData, err := os.ReadFile("server.crt")
	if err != nil {
		log.Fatal("ReadFile: ", err)

	}
	certs.AppendCertsFromPEM(pemData)
	mTLSConfig := tls.Config{
		RootCAs: certs,
	}

	fmt.Println("start 1", mTLSConfig)
	http.HandleFunc("/hello", HelloServer)
	// err = http.ListenAndServe(":443", nil)
	err = http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
