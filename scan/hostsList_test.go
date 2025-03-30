package scan_test

import (
	"errors"
	"os"
	"testing"

	"vegorov.ru/go-cli/pScan/scan"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаём пустой список хостов для теста
			hl := &scan.HostsList{}

			// Добавляем первый хост
			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}

			err := hl.Add(tc.host)

			// Спачала проверяем на ожидаемые ошибки
			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("Ожидали ошибку %q, а получили nil", tc.expectErr)
				}

				if !errors.Is(err, tc.expectErr) {
					t.Errorf("Ожидали ошибку %q, а получили %q\n", tc.expectErr, err)
				}

				// Завершаем проверку на ожидаемые ошибки
				return
			}

			// Если мы дошли до этого места - то ошибок в тесте мы не ожидали
			// поэтому нужно проверить, нет ли неожиданной ошибки
			if err != nil {
				t.Fatalf("Ошибок в тесте не ожидали, а получили ошибку: %q", err)
			}

			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Ожидали длину списка хостов %d, а получили %d\n", tc.expectLen, len(hl.Hosts))
			}

			if hl.Hosts[1] != tc.host {
				t.Errorf("Ожидали имя хоста %q по индексу 1, а получили %q\n", tc.host, hl.Hosts[1])
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		name        string
		host        string
		expectedLen int
		expectedErr error
	}{
		{"RemoveExisting", "host1", 1, nil},
		{"RemoveNotFound", "host3", 1, scan.ErrNotExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}

			for _, h := range []string{"host1", "host2"} {
				if err := hl.Add(h); err != nil {
					t.Fatal(err)
				}
			}

			err := hl.Remove(tc.host)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("Ожидали ошибку %q, а получили nil", tc.expectedErr)
				}

				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("Ожидали ошибку %q, а получили %q\n", tc.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Не ожидали получить ошибку, а получили %q\n", err)
			}

			if len(hl.Hosts) != tc.expectedLen {
				t.Errorf("Ожидали длину списка хостов %d, а получили %d\n", tc.expectedLen, len(hl.Hosts))
			}

			// es := "Имя хоста %i не должно быть в списке хостов, а оно там есть\n"
			if hl.Hosts[0] == tc.host {
				t.Errorf("Имя хоста %q не должно быть в списке хостов, а оно там есть\n", tc.host)
			}
		})
	}
}

func TestSaveLoad(t *testing.T) {
	hl1 := scan.HostsList{}
	hl2 := scan.HostsList{}

	hostName := "host1"
	hl1.Add(hostName)

	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %s\n", err)
	}
	defer os.Remove(tf.Name())

	if err := hl1.Save(tf.Name()); err != nil {
		t.Fatalf("Ошибка сохранения списка хостов в файл: %s", err)
	}

	if err := hl2.Load(tf.Name()); err != nil {
		t.Fatalf("Ошибка чтения списка хостов из файла: %s", err)
	}

	if hl1.Hosts[0] != hl2.Hosts[0] {
		t.Errorf("Имена хостов не совпадают: %q и %q", hl1.Hosts[0], hl2.Hosts[0])
	}
}

func TestLoadNoFile(t *testing.T) {
	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %s\n", err)
	}

	if err := os.Remove(tf.Name()); err != nil {
		t.Fatalf("Ошибка удаления временного файла: %s\n", err)
	}

	hl := &scan.HostsList{}

	if err := hl.Load(tf.Name()); err != nil {
		t.Errorf("Не ожидали ошибку, а получили %q\n", err)
	}
}
