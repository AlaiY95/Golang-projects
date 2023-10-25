package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// Create a scanner to read input from the console
	scanner := bufio.NewScanner(os.Stdin)
	// Print a header for the output
	fmt.Printf("domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord\n")

	// Read input line by line
	for scanner.Scan() {
		checkDomain(scanner.Text()) // Call the checkDomain function for each domain
	}

	// Check for errors in the input scanner
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	// Look up MX (Mail Exchanger) records for the domain
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	// If MX records exist, set hasMX to true
	if len(mxRecords) > 0 {
		hasMX = true
	}

	// Look up TXT records for the domain
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	// Iterate over TXT records
	for _, record := range txtRecords {
		// Check if a TXT record starts with "v=spf1" to determine the presence of SPF
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	// Look up DMARC records for the domain by querying a specific subdomain
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	// Iterate over DMARC records
	for _, record := range dmarcRecords {
		// Check if a DMARC record starts with "v=DMARC1" to determine the presence of DMARC
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	// Print the results for the domain, including presence flags and associated records
	fmt.Printf("%v, %v, %v, %v, %v, %v\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
