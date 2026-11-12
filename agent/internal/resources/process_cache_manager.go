package resources

import "sync"

type ProcessCacheData struct {
	Name     string
	Username string
	Command  string
	Ppid     int32
	Uptime   int64
}
type ProcessCacheManager struct {
	lock  sync.Mutex
	cache map[int32]ProcessCacheData
}

func NewProcessCacheManager() *ProcessCacheManager {
	return &ProcessCacheManager{
		cache: make(map[int32]ProcessCacheData),
	}
}

func (m *ProcessCacheManager) Get(pid int32) (ProcessCacheData, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	data, ok := m.cache[pid]
	return data, ok
}

func (m *ProcessCacheManager) Set(pid int32, data ProcessCacheData) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.cache[pid] = data
}

func (m *ProcessCacheManager) Delete(pid int32) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.cache, pid)
}

func (m *ProcessCacheManager) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.cache = make(map[int32]ProcessCacheData)
}

func (m *ProcessCacheManager) Keys() []int32 {
	m.lock.Lock()
	defer m.lock.Unlock()

	keys := make([]int32, 0, len(m.cache))
	for k := range m.cache {
		keys = append(keys, k)
	}
	return keys
}

func (m *ProcessCacheManager) RemoveAllBut(newKeys []int32) {

	newKeySet := make(map[int32]struct{})
	for _, k := range newKeys {
		newKeySet[k] = struct{}{}
	}
	for _, k := range m.Keys() {
		// if key is not in newKeys, delete it
		if _, ok := newKeySet[k]; !ok {
			m.Delete(k)
		}
	}
}
