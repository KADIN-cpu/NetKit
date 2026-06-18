package portscan

import "testing"

func TestPortRange(t *testing.T) {
	ports := PortRange(80, 83)
	want := []int{80, 81, 82, 83}
	if len(ports) != len(want) {
		t.Fatalf("esperava %d portas, obteve %d", len(want), len(ports))
	}
	for i, p := range ports {
		if p != want[i] {
			t.Errorf("ports[%d] = %d, esperava %d", i, p, want[i])
		}
	}
}

func TestPortRangeReversed(t *testing.T) {
	// Deve normalizar limites invertidos.
	ports := PortRange(83, 80)
	if len(ports) != 4 {
		t.Fatalf("esperava 4 portas, obteve %d", len(ports))
	}
}

func TestServiceName(t *testing.T) {
	if got := ServiceName(22); got != "ssh" {
		t.Errorf("porta 22 = %q, esperava \"ssh\"", got)
	}
	if got := ServiceName(6379); got != "redis" {
		t.Errorf("porta 6379 = %q, esperava \"redis\"", got)
	}
	if got := ServiceName(12345); got != "?" {
		t.Errorf("porta desconhecida = %q, esperava \"?\"", got)
	}
}
