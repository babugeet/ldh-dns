package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings" // Import the strings package
)

func main() {
	// Change the port number here
	addr := net.UDPAddr{
		Port: 8053,                   // Example: Use port 8053 instead of 53
		IP:   net.ParseIP("0.0.0.0"), // Listen on all available IPs
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Error setting up UDP listener: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Printf("DNS server listening on port %d...\n", addr.Port)

	buffer := make([]byte, 512) // Buffer to hold incoming DNS queries

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading UDP packet: %v\n", err)
			continue
		}

		fmt.Printf("Received DNS query from %v, %d bytes\n", clientAddr, n)

		// Parse the DNS query
		if err := parseDNSQuery(buffer[:n], conn, clientAddr); err != nil {
			fmt.Printf("Error parsing DNS query: %v\n", err)
			continue
		}
	}
}

// parseDNSQuery parses the incoming DNS query
func parseDNSQuery(data []byte, conn *net.UDPConn, clientAddr *net.UDPAddr) error {
	if len(data) < 12 {
		return fmt.Errorf("invalid DNS query length")
	}

	// 1. Parse DNS header (12 bytes)
	transactionID := binary.BigEndian.Uint16(data[0:2])
	flags := binary.BigEndian.Uint16(data[2:4])
	questions := binary.BigEndian.Uint16(data[4:6])

	fmt.Printf("Transaction ID: %d\n", transactionID)
	fmt.Printf("Flags: %d\n", flags)
	fmt.Printf("Questions: %d\n", questions)

	// 2. Parse the Question section
	offset := 12
	for i := 0; i < int(questions); i++ {
		// Extract the domain name
		domainName, newOffset, err := parseDomainName(data, offset)
		if err != nil {
			return err
		}
		offset = newOffset

		// Get the question type and class
		if len(data) < offset+4 {
			return fmt.Errorf("invalid question section")
		}
		qtype := binary.BigEndian.Uint16(data[offset : offset+2])
		qclass := binary.BigEndian.Uint16(data[offset+2 : offset+4])

		fmt.Printf("Domain Name: %s\n", domainName)
		fmt.Printf("Query Type: %d\n", qtype)
		fmt.Printf("Query Class: %d\n", qclass)

		offset += 4
		// Craft a DNS response
		response := createDNSResponse(transactionID, domainName)

		// Send the response to the client
		_, err = conn.WriteToUDP(response, clientAddr)
		if err != nil {
			return fmt.Errorf("error sending DNS response: %v", err)
		}
	}
	return nil
}

// createDNSResponse creates a DNS response based on the query
func createDNSResponse(transactionID uint16, domainName string) []byte {
	// DNS Header (12 bytes)
	header := make([]byte, 12)
	binary.BigEndian.PutUint16(header[0:2], transactionID) // Transaction ID
	binary.BigEndian.PutUint16(header[2:4], 0x8180)        // Flags: Standard query response, no error
	binary.BigEndian.PutUint16(header[4:6], 1)             // Number of questions
	binary.BigEndian.PutUint16(header[6:8], 1)             // Number of answers
	binary.BigEndian.PutUint16(header[8:10], 0)            // Number of authority RRs
	binary.BigEndian.PutUint16(header[10:12], 0)           // Number of additional RRs

	// Question Section
	question := createDomainNameSection(domainName)     // The domain name section
	question = append(question, 0x00, 0x01, 0x00, 0x01) // Type A, Class IN

	// Answer Section
	// The answer will be a simple static response. For example, returning an IP address "127.0.0.1"
	answer := createDomainNameSection(domainName)   // Domain Name
	answer = append(answer, 0x00, 0x01)             // Type A
	answer = append(answer, 0x00, 0x01)             // Class IN
	answer = append(answer, 0x00, 0x00, 0x00, 0x1e) // TTL (30 seconds)
	answer = append(answer, 0x00, 0x04)             // Data length (4 bytes for IPv4 address)
	answer = append(answer, 0x7f, 0x00, 0x00, 0x01) // IP address 127.0.0.1

	// Combine all parts into the response
	response := append(header, question...)
	response = append(response, answer...)

	return response
}

// createDomainNameSection creates the domain name section for the DNS message
func createDomainNameSection(domainName string) []byte {
	parts := strings.Split(domainName, ".")
	var section []byte
	for _, part := range parts {
		section = append(section, byte(len(part)))
		section = append(section, part...)
	}
	section = append(section, 0x00) // End of domain name
	return section
}

// parseDomainName extracts the domain name from the DNS question section
func parseDomainName(data []byte, offset int) (string, int, error) {
	var domainParts []string
	for {
		if offset >= len(data) {
			return "", 0, fmt.Errorf("invalid domain name")
		}
		length := int(data[offset])
		if length == 0 {
			offset++
			break
		}
		offset++
		if offset+length > len(data) {
			return "", 0, fmt.Errorf("invalid domain name length")
		}
		domainParts = append(domainParts, string(data[offset:offset+length]))
		offset += length
	}

	// Use strings.Join to combine the domain parts
	return strings.Join(domainParts, "."), offset, nil
}
