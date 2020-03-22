package botv2

import kitlog "github.com/go-kit/kit/log"

type Logger struct {
	kitlog.Logger
}

func (l *Logger) Message(s string) {
	l.Log("msg", s)
}
