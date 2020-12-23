
type logger struct {
	indent string
	Suppress bool
}
var log logger

func (l *logger) Push() {
	l.indent += "   "
}

func (l *logger) Pop() {
	l.indent = l.indent[:len(l.indent)-3]
}

func (l *logger) Println(args ...interface{}) {
	if l.Suppress {
		return
	}
	if len(l.indent) > 0 {
		fmt.Println(append([]interface{}{l.indent[:len(l.indent)-1]}, args...)...)
	} else {
		fmt.Println(args...)
	}
}

func (l *logger) Print(args ...interface{}) {
	if l.Suppress {
		return
	}
	if len(l.indent) > 0 {
		fmt.Print(append([]interface{}{l.indent[:len(l.indent)-1]}, args...)...)
	} else {
		fmt.Print(args...)
	}
}

func (l *logger) Printf(args ...interface{}) {
	if l.Suppress {
		return
	}
	if len(l.indent) > 0 {
		fmt.Printf("%s"+args[0].(string), append([]interface{}{l.indent}, args[1:]...)...)
	} else {
		fmt.Printf(args[0].(string), args[1:]...)
	}
}
