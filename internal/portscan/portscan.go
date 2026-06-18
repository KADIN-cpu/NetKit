package portscan

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

// PortResult representa o estado de uma porta.
type PortResult struct {
	Port    int
	Open    bool
	Service string
}

// commonServices mapeia portas conhecidas a nomes de serviço.
var commonServices = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	389:  "ldap",
	443:  "https",
	445:  "smb",
	636:  "ldaps",
	1433: "mssql",
	3000: "node/dev",
	3306: "mysql",
	3389: "rdp",
	5432: "postgres",
	5900: "vnc",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	9090: "prometheus",
}

// ServiceName retorna o nome do serviço de uma porta, ou "?" se desconhecido.
func ServiceName(port int) string {
	if s, ok := commonServices[port]; ok {
		return s
	}
	return "?"
}

// Scan testa um conjunto de portas em um host com concorrência limitada.
func Scan(host string, ports []int, timeout time.Duration, workers int) []PortResult {
	if workers <= 0 {
		workers = 100
	}

	results := make([]PortResult, len(ports))
	jobs := make(chan int, len(ports))
	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				port := ports[i]
				results[i] = PortResult{
					Port:    port,
					Open:    isOpen(host, port, timeout),
					Service: ServiceName(port),
				}
			}
		}()
	}

	for i := range ports {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	return results
}

func isOpen(host string, port int, timeout time.Duration) bool {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// CommonPorts retorna uma lista das portas mais usadas para um scan rápido.
func CommonPorts() []int {
	ports := make([]int, 0, len(commonServices))
	for p := range commonServices {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	return ports
}

// PortRange gera uma faixa de portas de start até end (inclusive).
func PortRange(start, end int) []int {
	if start > end {
		start, end = end, start
	}
	ports := make([]int, 0, end-start+1)
	for p := start; p <= end; p++ {
		if p >= 1 && p <= 65535 {
			ports = append(ports, p)
		}
	}
	return ports
}
