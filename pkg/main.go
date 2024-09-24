package main

import (
	"fmt"
	"ldh-dns/pkg/internal/dns"
	"net"
	// Import the strings package
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
		if err := dns.ParseDNSQuery(buffer[:n], conn, clientAddr); err != nil {
			fmt.Printf("Error parsing DNS query: %v\n", err)
			continue
		}
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

// // resolveToLinuxDataHub resolves linuxdatahub.local to linuxdatahub.com
// func resolveToLinuxDataHub(domain string) string {
// 	fmt.Printf("Resolving %s to linuxdatahub.com\n", domain)
// 	// Use the public IP or any IP associated with linuxdatahub.com
// 	// This would be the internal resolution
// 	return "ldh.svc" // Example IP for linuxdatahub.com (replace with actual)
// }

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
