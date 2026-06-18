# netkit

> Toolkit de infraestrutura e rede em linha de comando — um único binário, multiplataforma.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build](https://github.com/KADIN-cpu/netkit/actions/workflows/build.yml/badge.svg)](https://github.com/KADIN-cpu/netkit/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/KADIN-cpu/netkit)](https://goreportcard.com/report/github.com/KADIN-cpu/netkit)

`netkit` reúne as tarefas de diagnóstico de rede e sistema mais comuns num único binário, sem depender de scripts soltos espalhados pela máquina. Funciona em **Windows, Linux e macOS**.

## Recursos

- **`sysinfo`** — informações completas do sistema: CPU, memória, discos e uptime
- **`ping`** — ping em um host único ou varredura de uma faixa de rede inteira via CIDR (concorrente)
- **`portscan`** — scanner de portas TCP com identificação de serviços comuns

## Motivação

Comecei este projeto como uma forma de aprender Go na prática construindo algo que eu realmente usaria. No dia a dia, sempre acabava alternando entre vários utilitários diferentes (`ping`, `nmap`, `systeminfo`, `df -h`) e alguns scripts próprios só para checar o básico de uma máquina ou de uma rede.

A ideia do `netkit` foi juntar o essencial num comando só, com saída limpa e colorida, num binário único que roda igual no Windows e no Linux — sem instalar runtime nenhum.

## Instalação

### Binário pré-compilado

Baixe o binário para o seu sistema na [página de releases](https://github.com/KADIN-cpu/netkit/releases) e coloque no seu `PATH`.

### Via `go install`

```bash
go install github.com/KADIN-cpu/netkit@latest
```

### Compilando do código-fonte

```bash
git clone https://github.com/KADIN-cpu/netkit.git
cd netkit
make build      # ou: go build -o netkit .
```

## Uso

### Informações do sistema

```bash
netkit sysinfo
```

```
» Sistema
  Hostname      meu-pc
  SO            linux
  Plataforma    ubuntu 22.04
  Arquitetura   x86_64
  Uptime        12d 4h 33m

» CPU
  Modelo        Intel(R) Core(TM) i7-10700
  Núcleos       8
  Uso           14.2%

» Memória
  Total         16.0 GiB
  Em uso        9.4 GiB
  Uso           58.7%
```

### Ping / varredura de rede

```bash
# host único
netkit ping 192.168.0.10

# faixa inteira (CIDR)
netkit ping 192.168.0.0/24

# rede /24, mostrando só quem está online
netkit ping 192.168.0.0/24 --only-up
```

Flags: `--timeout/-t` (ms por host), `--workers/-w` (concorrência), `--only-up`.

### Scanner de portas

```bash
# portas comuns
netkit portscan 192.168.0.10

# portas específicas
netkit portscan 192.168.0.10 --ports 22,80,443

# faixa de portas
netkit portscan 192.168.0.10 --ports 1-1024

# tudo (1-65535)
netkit portscan 192.168.0.10 --all
```

## Build multiplataforma

```bash
make build-all     # gera binários para windows, linux e macos em ./dist
```

Ou manualmente:

```bash
GOOS=windows GOARCH=amd64 go build -o dist/netkit.exe .
GOOS=linux   GOARCH=amd64 go build -o dist/netkit-linux .
GOOS=darwin  GOARCH=arm64 go build -o dist/netkit-macos .
```

## Stack

- [Go 1.22](https://go.dev)
- [Cobra](https://github.com/spf13/cobra) — estrutura de comandos CLI
- [gopsutil](https://github.com/shirou/gopsutil) — métricas de sistema multiplataforma
- [tablewriter](https://github.com/olekukonko/tablewriter) — tabelas no terminal
- [color](https://github.com/fatih/color) — saída colorida

## Roadmap

- [ ] Exportação dos resultados em JSON / CSV
- [ ] Comando `traceroute`
- [ ] Resolução reversa de DNS na varredura
- [ ] Modo de monitoramento contínuo (`watch`)

## Licença

MIT — veja [LICENSE](LICENSE).
