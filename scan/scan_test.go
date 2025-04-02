package scan_test

import (
	"net"
	"strconv"
	"testing"

	"vegorov.ru/go-cli/pScan/scan"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}

	if ps.Open.String() != "closed" {
		t.Errorf("Ожидали: %q, получили: %q/n", "closed", ps.Open.String())
	}

	ps.Open = true
	if ps.Open.String() != "open" {
		t.Errorf("Ожидали: %q, получили: %q/n", "open", ps.Open.String())
	}
}

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name          string
		expectedState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}

	host := "localhost"
	hl := &scan.HostsList{}
	hl.Add(host)

	ports := []int{}

	for _, tc := range testCases {
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)

		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}

	res := scan.Run(hl, ports)

	if len(res) != 1 {
		t.Fatalf("Ожидали 1 результат, а получили: %d\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("Ожидали хост: %q, а получили: %q\n", host, res[0].Host)
	}

	if res[0].NotFound {
		t.Errorf("Ожидаемый хост %q не найден\n", host)
	}

	if len(res[0].PortStates) != 2 {
		t.Fatalf("Ожидали 2 состояния портов, получили %d\n", len(res[0].PortStates))
	}

	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Ожидали порт: %d, получили: %d\n", ports[i], res[0].PortStates[i].Port)
		}

		if res[0].PortStates[i].Open.String() != tc.expectedState {
			t.Errorf("Ожидали порт %d в состоянии: %s\n", ports[i], tc.expectedState)
		}
	}
}

func TestRunHostNotFound(t *testing.T) {
	host := "257.257.257.257"
	hl := &scan.HostsList{}
	hl.Add(host)

	res := scan.Run(hl, []int{})

	if len(res) != 1 {
		t.Fatalf("Ожидали 1 результат, получили: %d\n", len(res))
	}

	if res[0].Host != host {
		t.Fatalf("Ожидали хост: %q, получили: %q\n", host, res[0].Host)
	}

	if !res[0].NotFound {
		t.Errorf("Ожидали, что хост %q будет не найден\n", host)
	}

	if len(res[0].PortStates) != 0 {
		t.Fatalf("Ожидали 0 состояний портов, получили: %d\n", len(res[0].PortStates))
	}
}
