package scan

import (
	"errors"
	"sort"
)

var (
	ErrExists    = errors.New("Хост уже в списке")
	ErrNotExists = errors.New("Хост не в списке")
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
