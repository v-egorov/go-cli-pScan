package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
)

var (
	ErrExists    = errors.New("хост уже в списке")
	ErrNotExists = errors.New("хост не в списке")
)

// Список хостов для сканирования
type HostsList struct {
	Hosts []string
}

// search выполняет поиск в списке хостов
func (hl *HostsList) search(host string) (bool, int) {
	sort.Strings(hl.Hosts)

	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}
	return false, -1
}

// Add добавляет хост в список
func (hl *HostsList) Add(host string) error {
	if found, _ := hl.search(host); found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}
	hl.Hosts = append(hl.Hosts, host)
	return nil
}

func (hl *HostsList) Remove(host string) error {
	if found, i := hl.search(host); found {
		hl.Hosts = slices.Delete(hl.Hosts, i, i+1)
		return nil
	}
	return fmt.Errorf("%w: %s", ErrNotExists, host)
}

// Load загружает список хостов из файла
func (hl *HostsList) Load(hostsFile string) error {
	f, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}

	return nil
}

// Save сохраняет список хостов в файл
func (hl *HostsList) Save(hostsFile string) error {
	output := ""
	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}

	return os.WriteFile(hostsFile, []byte(output), 0644)
}
