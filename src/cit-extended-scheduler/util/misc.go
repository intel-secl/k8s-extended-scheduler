package util

import (
	"flag"
	"strconv"
)

var AH_KEY_FILE string

func GetCmdlineArgs() (*string, string, *string, *string) {

	// IP and PORT for the extended scheduler to listen.
	url := flag.String("url", "127.0.0.1", "IP address for the extended scheduler to listen on")
	port_no := flag.Int("port", 8888, "Port number for the extended scheduler to listen on")
	server_crt := flag.String("server_crt", "", "Server Certificate to be used for TLS handshake ")
	server_key := flag.String("server_key", "", "Server Key to be used for TLS handshake ")
	ah_key := flag.String("ah_key", "", "Attestation Hub Key to be used for parsing signed trust report")

	//parse the cmdline args
	flag.Parse()
	port := strconv.Itoa(*port_no)

	AH_KEY_FILE = (*ah_key)
	return url, port, server_crt, server_key
}
