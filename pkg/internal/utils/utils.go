package utils

import (
	"fmt"
	"strings"
)

func SplitDomain(input string) string {
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

// Instead of resolving to an IP, resolve the query to the CNAME linuxdatahub.svc
func ResolveToLinuxDataHub(domain string) string {
	localsvc := SplitDomain(domain)
	fmt.Printf("Resolving %s to CNAME linuxdatahub.svc\n", domain)
	fmt.Println(localsvc)
	// Return the CNAME for linuxdatahub.svc
	return localsvc
}

// createDomainNameSection creates the domain name section for the DNS message
func CreateDomainNameSection(domainName string) []byte {
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
func ParseDomainName(data []byte, offset int) (string, int, error) {
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
