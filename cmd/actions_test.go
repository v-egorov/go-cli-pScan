package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

	// Интеграционный тест: add --> list --> delete --> list
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

	if out.String() != expectedOut {
		t.Errorf("Ожидали вывод: %q, получили: %q\n", expectedOut, out.String())
	}
}
