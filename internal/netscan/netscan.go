package netscan

import (
	"fmt"
	"net"
)

// HostsFromCIDR expande uma notação CIDR (ex: 192.168.0.0/24) em
// todos os IPs de host utilizáveis, excluindo rede e broadcast.
func HostsFromCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("CIDR inválido: %w", err)
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// Para redes /31 e /32 retorna como está; senão remove rede e broadcast.
	if len(ips) <= 2 {
		return ips, nil
	}
	return ips[1 : len(ips)-1], nil
}

// inc incrementa um endereço IP em 1.
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
