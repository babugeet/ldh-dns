// package main

// import (
// 	"fmt"
// 	"strings"
// )

// func main() {
// 	input := "test.ns.linuxdatahub.local"
// 	// input := "test.linuxdatahub.local"
// 	//test.linuxdatahub.svc.cluster.local.
// 	fmt.Println(input)
// 	parts := strings.Split(input, ".")
// 	if len(input) > 4 {
// 		part := parts[:2]
// 		fmt.Println(strings.Join(part, ".") + ".svc.cluster.local")
// 	} else if len(input) == 3 {
// 		part := parts[:1]
// 		fmt.Println(strings.Join(part, ".") + "linuxdatahub.svc.cluster.local")
// 	}
// 	fmt.Println(parts)

// }
