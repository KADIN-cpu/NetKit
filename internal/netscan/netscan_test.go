package netscan

import "testing"

func TestHostsFromCIDR(t *testing.T) {
	hosts, err := HostsFromCIDR("192.168.0.0/30")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	// /30 = 4 endereços; removendo rede e broadcast sobram 2.
	want := []string{"192.168.0.1", "192.168.0.2"}
	if len(hosts) != len(want) {
		t.Fatalf("esperava %d hosts, obteve %d (%v)", len(want), len(hosts), hosts)
	}
	for i, h := range hosts {
		if h != want[i] {
			t.Errorf("host[%d] = %q, esperava %q", i, h, want[i])
		}
	}
}

func TestHostsFromCIDRInvalid(t *testing.T) {
	if _, err := HostsFromCIDR("not-a-cidr"); err == nil {
		t.Error("esperava erro para CIDR inválido, obteve nil")
	}
}
