package ping

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// Result representa o resultado de um ping a um host.
type Result struct {
	Host  string
	Alive bool
	RTT   time.Duration
}

// Ping executa um único ping ao host usando o utilitário nativo do SO.
// Isso evita a necessidade de privilégios de root para sockets ICMP crus.
func Ping(ctx context.Context, host string, timeout time.Duration) Result {
	res := Result{Host: host}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// -n 1 (1 pacote), -w timeout em ms
		ms := fmt.Sprintf("%d", timeout.Milliseconds())
		cmd = exec.CommandContext(ctx, "ping", "-n", "1", "-w", ms, host)
	default:
		// -c 1 (1 pacote), -W timeout em segundos (mínimo 1)
		secs := int(timeout.Seconds())
		if secs < 1 {
			secs = 1
		}
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", "-W", fmt.Sprintf("%d", secs), host)
	}

	start := time.Now()
	err := cmd.Run()
	res.RTT = time.Since(start)
	res.Alive = err == nil
	return res
}

// Sweep executa ping concorrente em vários hosts, limitando a concorrência.
func Sweep(hosts []string, timeout time.Duration, workers int) []Result {
	if workers <= 0 {
		workers = 50
	}

	results := make([]Result, len(hosts))
	jobs := make(chan int, len(hosts))
	var wg sync.WaitGroup

	ctx := context.Background()

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				results[i] = Ping(ctx, hosts[i], timeout)
			}
		}()
	}

	for i := range hosts {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	return results
}
