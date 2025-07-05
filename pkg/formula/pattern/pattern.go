package pattern

const (
	Number   string = `([0-9]+(\.[0-9]+)?(e[\+-][0-9]+)?)`
	Function string = `([a-zA-Z][a-zA-Z0-9]*)`
	Variable string = `\[([a-zA-Z0-9-:]+)\]`
)

func Dict() []string {
	return []string{Number, Function, Variable}
}
