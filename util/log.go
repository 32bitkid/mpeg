package util

import deflog "log"
import "io/ioutil"
import "os"
import "strings"

var namespaces = strings.Split(os.Getenv("DEBUG"), ",")

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

var NoopLogger = deflog.New(ioutil.Discard, "", 0)

func NewLog(ns string) Logger {
	for _, namespace := range namespaces {
		if namespace == "*" {
			return createLogger(ns)
		}

		if strings.HasSuffix(namespace, ":*") && strings.HasPrefix(ns, namespace[:len(namespace)-2]) {
			return createLogger(ns)
		}

		if ns == namespace {
			return createLogger(ns)
		}
	}
	return NoopLogger
}

func createLogger(ns string) Logger {
	return deflog.New(os.Stdout, ns+": ", 0)
}
