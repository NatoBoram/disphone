package main

// rsfa : Remove String From Array.
func rsfa(a []string, s string) []string {
	var n []string
	for i, v := range a {
		if v != s {
			n = append(n, a[i])
		}
	}
	return n
}

// csia : Check String Inside Array
func csia(a []string, s string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}
	return false
}
