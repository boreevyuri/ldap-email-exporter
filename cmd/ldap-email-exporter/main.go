package main

import (
	"flag"
	"log"
	"os"

	"ldap-email-exporter/cmd/ldap-email-exporter/config"

	"ldap-email-exporter/cmd/ldap-email-exporter/ldapsearch"
)

var configFile = flag.String("c", "./email-exporter.yaml", "-c ./email-exporter.yaml")

func main() {
	flag.Parse()
	config := configuration.New(*configFile)
	errLog := log.New(os.Stderr, "#ERROR: ", log.Flags())

	ls, err := ldapsearch.New(config.LDAP)
	if err != nil {
		// Print error and exit
		errLog.Printf("Unable to connect to LDAP server: %s", err)
		ls.Close()
		os.Exit(1)
	}
	defer ls.Close()

	err = ls.Search()
	if err != nil {
		// Print error and exit
		errLog.Printf("Got error in search params: %s", err)
		ls.Close()
		os.Exit(1)
	}

	ls.Print()
}
