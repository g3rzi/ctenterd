package builtin

import "sort"

type Fn func([]string) //int

type entry struct {
	Run   Fn
	Help  string
	Alias []string
}

var reg = map[string]entry{}

func Register(name string, fn Fn, help string, aliases ...string) {
	reg[name] = entry{Run: fn, Help: help, Alias: aliases}
	for _, a := range aliases {
		reg[a] = entry{Run: fn, Help: help + " (alias: " + name + ")", Alias: nil}
	}
}

func Lookup(name string) (Fn, bool) {
	e, ok := reg[name]
	return e.Run, ok
}

func List() (names []string) {
	seen := map[string]bool{}
	for k, e := range reg {
		// hide alias keys when their help marks them as alias
		if len(e.Alias) == 0 && !seen[k] {
			names = append(names, k)
			seen[k] = true
		}
	}
	sort.Strings(names)
	return names
}

func Help(name string) (string, bool) {
	e, ok := reg[name]
	return e.Help, ok
}
