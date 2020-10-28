package utils

func UniqueStr(strs []string) []string {
	m := make(map[string]struct{})
	for _, s := range strs {
		m[s] = struct{}{}
	}
	res := make([]string, 0)
	for k := range m {
		res = append(res, k)
	}
	return res
}

func FilterStr(strs []string, f func(s string) bool) (res []string) {
	for _, s := range strs {
		if f(s) {
			res = append(res, s)
		}
	}
	return
}
