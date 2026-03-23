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
		// Timeout court pour ne pas attendre indéfiniment
		conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)

		if err == nil {
			fmt.Printf("[+] Port %d ouvert\n", port)
			conn.Close()
		}
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <host> <port_debut> <port_fin>")
		fmt.Println("Exemple: go run main.go 127.0.0.1 22 100")
		return
	}

	host := os.Args[1]
	start, _ := strconv.Atoi(os.Args[2])
	end, _ := strconv.Atoi(os.Args[3])

	if start > end {
		fmt.Println("Erreur: Le port de début doit être inférieur au port de fin.")
		return
	}

	// Calcul de la division pour le double threading
	totalPorts := end - start + 1
	midPoint := start + (totalPorts / 2) - 1

	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Printf("Scan de %s (Ports %d à %d) avec 2 threads...\n", host, start, end)
	startWait := time.Now()

	// Thread A : Première moitié
	go scanPorts(host, start, midPoint, &wg)

	// Thread B : Deuxième moitié
	go scanPorts(host, midPoint+1, end, &wg)

	wg.Wait()
	fmt.Printf("\nScan terminé en %v\n", time.Since(startWait))
}
