package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings" // Import the strings package
	"time"
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
		// 	// Craft a DNS response
		// 	response := createDNSResponse(transactionID, domainName)

		// 	// Send the response to the client
		// 	_, err = conn.WriteToUDP(response, clientAddr)
		// 	if err != nil {
		// 		return fmt.Errorf("error sending DNS response: %v", err)
		// 	}
		// }
		// return nil
		// If the domain is within linuxdatahub.local, resolve it locally
		if strings.HasSuffix(domainName, "sample.dev.ldhappdomain.cloud") || domainName == "sample.dev.ldhappdomain.cloud:" {
			cname := resolveToLinuxDataHub(domainName)

			response := createDNSResponse(transactionID, domainName, cname)
			_, err := conn.WriteToUDP(response, clientAddr)
			if err != nil {
				return fmt.Errorf("error sending DNS response: %v", err)
			}
		} else {
			// Forward the request to another DNS server (e.g., 8.8.8.8)
			forwardToExternalDNS(data, conn, clientAddr)
		}
	}
	return nil
}

// forwardToExternalDNS forwards a DNS query to another DNS server (e.g., Google DNS)
func forwardToExternalDNS(query []byte, conn *net.UDPConn, clientAddr *net.UDPAddr) {
	dnsServerAddr := net.UDPAddr{
		IP:   net.ParseIP("10.96.0.10"), // Forward to Google DNS (or other DNS server)
		Port: 53,
	}

	// Create a connection to the external DNS server
	externalConn, err := net.DialUDP("udp", nil, &dnsServerAddr)
	if err != nil {
		fmt.Printf("Error connecting to external DNS server: %v\n", err)
		return
	}
	defer externalConn.Close()

	// Send the DNS query to the external DNS server
	_, err = externalConn.Write(query)
	if err != nil {
		fmt.Printf("Error forwarding DNS query: %v\n", err)
		return
	}

	// Wait for the response from the external DNS server
	response := make([]byte, 512)
	externalConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, _, err := externalConn.ReadFromUDP(response)
	if err != nil {
		fmt.Printf("Error reading response from external DNS server: %v\n", err)
		return
	}

	// Send the response back to the client
	_, err = conn.WriteToUDP(response[:n], clientAddr)
	if err != nil {
		fmt.Printf("Error sending response to client: %v\n", err)
		return
	}
}

// func splitDomain(domain string) string {
// 	parts := strings.Split(domain, ".")
// 	if len(parts) < 2 {
// 		return domain + ".linuxdatahub.svc.cluster.local"
// 	}

// 	// Join the last two parts for the first result
// 	first := strings.Join(parts[len(parts)-2:], ".")

// 	// Join all parts except the la st two for the second result
// 	second := strings.Join(parts[:len(parts)-2], ",")

// 	return second + "." + first + ".svc.cluster.local"
// }

func splitDomain(input string) string {
	// input := "test.ns.linuxdatahub.local"
	// input := "test.linuxdatahub.local"
	//test.linuxdatahub.svc.cluster.local.
	fmt.Println(input)
	parts := strings.Split(input, ".")
	if len(input) > 4 {
		part := parts[:2]
		fmt.Println(strings.Join(part, ".") + ".svc.cluster.local")
		return strings.Join(part, ".") + ".svc.cluster.local"
	} else if len(input) == 3 {
		part := parts[:1]
		fmt.Println(strings.Join(part, ".") + "linuxdatahub.svc.cluster.local")
		return strings.Join(part, ".") + "linuxdatahub.svc.cluster.local"
	} else {
		// fmt.Println(parts)
		return input
	}
}

// // resolveToLinuxDataHub resolves linuxdatahub.local to linuxdatahub.com
// func resolveToLinuxDataHub(domain string) string {
// 	fmt.Printf("Resolving %s to linuxdatahub.com\n", domain)
// 	// Use the public IP or any IP associated with linuxdatahub.com
// 	// This would be the internal resolution
// 	return "ldh.svc" // Example IP for linuxdatahub.com (replace with actual)
// }

// Instead of resolving to an IP, resolve the query to the CNAME linuxdatahub.svc
func resolveToLinuxDataHub(domain string) string {
	localsvc := splitDomain(domain)
	fmt.Printf("Resolving %s to CNAME linuxdatahub.svc\n", domain)
	fmt.Println(localsvc)
	// Return the CNAME for linuxdatahub.svc
	return localsvc
}

// func createDNSResponse(transactionID uint16, domainName string, cname string) []byte {
// 	// DNS Header (12 bytes)
// 	header := make([]byte, 12)
// 	binary.BigEndian.PutUint16(header[0:2], transactionID) // Transaction ID
// 	binary.BigEndian.PutUint16(header[2:4], 0x8180)        // Flags: Standard query response, no error
// 	binary.BigEndian.PutUint16(header[4:6], 1)             // Number of questions
// 	binary.BigEndian.PutUint16(header[6:8], 1)             // Number of answers
// 	binary.BigEndian.PutUint16(header[8:10], 0)            // Number of authority RRs
// 	binary.BigEndian.PutUint16(header[10:12], 0)           // Number of additional RRs

// 	// Question Section
// 	question := createDomainNameSection(domainName)     // The domain name section
// 	question = append(question, 0x00, 0x01, 0x00, 0x01) // Type A, Class IN

// 	// Answer Section for CNAME
// 	answer := createDomainNameSection(domainName)   // Domain Name (the queried domain)
// 	answer = append(answer, 0x00, 0x05)             // Type CNAME (0x05)
// 	answer = append(answer, 0x00, 0x01)             // Class IN (0x01)
// 	answer = append(answer, 0x00, 0x00, 0x00, 0x1e) // TTL (30 seconds)

// 	cnameSection := createDomainNameSection(cname)                 // CNAME field
// 	dataLength := len(cnameSection)                                // Calculate correct data length
// 	answer = append(answer, byte(dataLength>>8), byte(dataLength)) // 2-byte length field for CNAME
// 	answer = append(answer, cnameSection...)                       // Append the CNAME data

// 	// Combine all parts into the response
// 	response := append(header, question...)
// 	response = append(response, answer...)

//		return response
//	}
//
// Modify the createDNSResponse function to return CNAME but keep .local in output
func createDNSResponse(transactionID uint16, domainName string, cname string) []byte {
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

	// Answer Section for CNAME
	answer := createDomainNameSection(domainName)   // Return domainName (.local) in the answer section
	answer = append(answer, 0x00, 0x05)             // Type CNAME (0x05)
	answer = append(answer, 0x00, 0x01)             // Class IN (0x01)
	answer = append(answer, 0x00, 0x00, 0x00, 0x1e) // TTL (30 seconds)

	cnameSection := createDomainNameSection(cname)                 // CNAME field (points to .svc)
	dataLength := len(cnameSection)                                // Calculate correct data length
	answer = append(answer, byte(dataLength>>8), byte(dataLength)) // 2-byte length field for CNAME
	answer = append(answer, cnameSection...)                       // Append the CNAME data

	// Combine all parts into the response
	response := append(header, question...)
	response = append(response, answer...)

	return response
}

// // createDNSResponse creates a DNS response for the given domain and IP address
// func createDNSResponse(transactionID uint16, domainName string, cname string) []byte {
// 	// DNS Header (12 bytes)
// 	header := make([]byte, 12)
// 	binary.BigEndian.PutUint16(header[0:2], transactionID) // Transaction ID
// 	binary.BigEndian.PutUint16(header[2:4], 0x8180)        // Flags: Standard query response, no error
// 	binary.BigEndian.PutUint16(header[4:6], 1)             // Number of questions
// 	binary.BigEndian.PutUint16(header[6:8], 1)             // Number of answers
// 	binary.BigEndian.PutUint16(header[8:10], 0)            // Number of authority RRs
// 	binary.BigEndian.PutUint16(header[10:12], 0)           // Number of additional RRs

// 	// Question Section
// 	question := createDomainNameSection(domainName)     // The domain name section
// 	question = append(question, 0x00, 0x01, 0x00, 0x01) // Type A, Class IN

// 	// Answer Section
// 	answer := createDomainNameSection(domainName)   // Domain Name
// 	answer = append(answer, 0x00, 0x05)             // Type CNAME (0x05)
// 	answer = append(answer, 0x00, 0x01)             // Class IN (0x01)
// 	answer = append(answer, 0x00, 0x00, 0x00, 0x1e) // TTL (30 seconds)
// 	// answer = append(answer, 0x00, 0x04)             // Data length (4 bytes for IPv4 address)
// 	// answer = append(answer, ip.To4()...)            // IP address (resolved or 127.0.0.1)
// 	cnameSection := createDomainNameSection(cname)   // CNAME field
// 	answer = append(answer, byte(len(cnameSection))) // Data length
// 	answer = append(answer, cnameSection...)

// 	// Combine all parts into the response
// 	response := append(header, question...)
// 	response = append(response, answer...)
// 	fmt.Println(response)
// 	return response
// }

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
