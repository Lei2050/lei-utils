package types

type MapMap[K1, K2 comparable, V any] struct {
	db map[K1]map[K2]V
}

func NewMapMap[K1, K2 comparable, V any]() *MapMap[K1, K2, V] {
	return &MapMap[K1, K2, V]{
		db: make(map[K1]map[K2]V),
	}
}

func (m *MapMap[K1, K2, V]) Push(key K1, nextKey K2, value V) {
	_, ok := m.db[key]
	if !ok {
		m.db[key] = make(map[K2]V)
	}
	m.db[key][nextKey] = value
}

func (m *MapMap[K1, K2, V]) Exit(key K1, nextKey K2) bool {
	gids, ok := m.db[key]
	if !ok {
		return false
	}
	_, ok = gids[nextKey]
	return ok
}

func (m *MapMap[K1, K2, V]) ExitFirstKey(key K1) bool {
	_, ok := m.db[key]
	return ok
}

func (m *MapMap[K1, K2, V]) GetValue(key K1, nextKey K2) V {
	gids, ok := m.db[key]
	if !ok {
		var v V
		return v
	}
	return gids[nextKey]
}

func (m *MapMap[K1, K2, V]) DelValue(key K1, nextKey K2) {
	_, ok := m.db[key]
	if !ok {
		return
	}
	delete(m.db[key], nextKey)
	if len(m.db[key]) == 0 {
		delete(m.db, key)
	}
}

func (m *MapMap[K1, K2, V]) DelFirstKey(key K1) {
	delete(m.db, key)
}

func (m *MapMap[K1, K2, V]) PrintValue(gid K1) map[K2]V {
	dbs, ok := m.db[gid]
	if !ok {
		return make(map[K2]V)
	}

	return dbs
}
func (m *MapMap[K1, K2, V]) Clear() {
	m.db = make(map[K1]map[K2]V)
}

func (m *MapMap[K1, K2, V]) GetFirstKeyValue(gid K1) map[K2]V {
	users, ok := m.db[gid]
	if !ok {
		return nil
	}
	return users
}

func (m *MapMap[K1, K2, V]) Values() map[K1]map[K2]V {
	return m.db
}
