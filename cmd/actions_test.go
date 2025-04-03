package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"vegorov.ru/go-cli/pScan/scan"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	// Создадим временный файл
	tf, err := os.CreateTemp("", "pScan")
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()

	if initList {
		hl := &scan.HostsList{}
		for _, h := range hosts {
			hl.Add(h)
		}
		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}
	}

	// Возвращаем имя временного файла и функцию его (этого файла) удаления
	// для очистки после выполнения теста
	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestHostActions(t *testing.T) {
	hosts := []string{
		"host1", "host2", "host3",
	}

	testCases := []struct {
		name           string
		args           []string
		expectedOut    string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOut:    "Добавлен хост: host1\nДобавлен хост: host2\nДобавлен хост: host3\n",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "ListAction",
			args:           hosts,
			expectedOut:    "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "DeleteAction",
			args:           []string{"host1", "host2"},
			expectedOut:    "Удалён хост: host1\nУдалён хост: host2\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tf, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()

			var out bytes.Buffer

			if err := tc.actionFunction(&out, tf, tc.args); err != nil {
				t.Fatalf("Не ожидали ошибку, а получили: %q\n", err)
			}

			if out.String() != tc.expectedOut {
				t.Errorf("Ожидали получить вывод: %q, а получили: %q", tc.expectedOut, out.String())
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	hosts := []string{
		"host1", "host2", "host3",
	}
	tf, cleanup := setup(t, hosts, false)
	defer cleanup()

	delHost := "host2"
	hostsEnd := []string{
		"host1", "host3",
	}

	var out bytes.Buffer

	expectedOut := ""
	for _, v := range hosts {
		expectedOut += fmt.Sprintf("Добавлен хост: %s\n", v)
	}

	// Формируем ожидаемый вывод всех операций интеграционного теста
	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintf("Удалён хост: %s\n", delHost)
	expectedOut += strings.Join(hostsEnd, "\n")
	expectedOut += fmt.Sprintln()

	for _, v := range hostsEnd {
		expectedOut += fmt.Sprintf("%s\nХост не найден", v)
		expectedOut += fmt.Sprintln()
	}

	// Интеграционный тест: add --> list --> delete --> scan --> list
	//
	// Add
	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatalf("Не ожидали ошибку, а получили: %q\n", err)
	}

	// List
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Не ожидали ошибку, а получили: %q\n", err)
	}

	// Delete host2
	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatalf("Не ожидали ошибку, а получили: %q\n", err)
	}

	// List после delete
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Не ожидали ошибку, а получили: %q\n", err)
	}

	// Сканируем хосты
	if err := scanAction(&out, tf, nil); err != nil {
		t.Fatalf("Не ожидали ошибку, а получили: %q\n", err)
	}

	// Проверим итоговый вывод
	if out.String() != expectedOut {
		t.Errorf("Ожидали вывод:\n%q, получили:\n%q\n", expectedOut, out.String())
	}
}

func TestScanAction(t *testing.T) {
	hosts := []string{
		"localhost",
		"not-found-host",
	}

	tf, cleanup := setup(t, hosts, true)
	defer cleanup()

	ports := []int{}
	// Подготовми порты, 1 открытый, 1 закрытый
	for i := 0; i < 2; i++ {
		ln, err := net.Listen("tcp", net.JoinHostPort("localhost", "0"))
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

		if i == 1 {
			ln.Close()
		}
	}

	expectedOut := fmt.Sprintln("localhost")
	expectedOut += fmt.Sprintf("\t%d: open\n", ports[0])
	expectedOut += fmt.Sprintf("\t%d: closed\n", ports[1])
	expectedOut += fmt.Sprintln("not-found-host\nХост не найден")

	var out bytes.Buffer

	if err := scanAction(&out, tf, ports); err != nil {
		t.Fatalf("Не ожидали ошибку, а получили: %q\n", err)
	}

	if out.String() != expectedOut {
		t.Errorf("Ожидали получить:\n%q, а получили:\n%q\n", expectedOut, out.String())
	}
}
