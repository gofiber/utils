package utils

const fastStrMapSliceLimit = 10000

type fastStrMap struct {
	data [fastStrMapSliceLimit]string
}

func NewFastMap(m map[string]string) *fastStrMap {
	fm := &fastStrMap{}
	fm.setMap(m)
	return fm
}

func (fm *fastStrMap) setMap(m map[string]string) {
	for k, v := range m {
		fm.set(k, v)
	}
}

func (fm *fastStrMap) set(key, val string) {
	fm.data[strToInt(key)] = val
}

func (fm *fastStrMap) Get(key string) string {
	i := strToInt(key)
	if i > fastStrMapSliceLimit {
		return ""
	}
	return fm.data[i]
}

func strToInt(s string) int {
	seed, hash := 131, 0

	for _, b := range s {
		hash = hash*seed + int(b)
	}
	return (hash & 0x7fffffff) % fastStrMapSliceLimit
}
