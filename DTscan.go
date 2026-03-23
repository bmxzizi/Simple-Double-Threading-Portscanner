package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func scanPorts(host string, startPort, endPort int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)

		if err == nil {
			fmt.Printf("[+] Port %d open.\n", port)
			conn.Close()
		}
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <host> <port_start> <port_end>")
		fmt.Println("Exemple: go run main.go 127.0.0.1 22 100")
		return
	}

	host := os.Args[1]
	start, _ := strconv.Atoi(os.Args[2])
	end, _ := strconv.Atoi(os.Args[3])

	if start > end {
		fmt.Println("Error: The start port must be inferior to the end port.")
		return
	}

	totalPorts := end - start + 1
	midPoint := start + (totalPorts / 2) - 1

	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Printf("Scanning %s (Ports %d to %d) ...\n", host, start, end)
	startWait := time.Now()

	go scanPorts(host, start, midPoint, &wg)

	go scanPorts(host, midPoint+1, end, &wg)

	wg.Wait()
	fmt.Printf("\nScan ended. en %v\n", time.Since(startWait))
}
