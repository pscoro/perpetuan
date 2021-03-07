package main

var commandMap = make(map[string]func(params []string) (bool, string))

func init() {
	commandMap["goto"] = func(params []string) (bool, string) {
		return false, ""
	}
}
