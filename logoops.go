package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	defaultFacility    = 20
	defaultHost        = "localhost"
	defaultMessage     = "hello"
	defaultPort        = 514
	defaultSeverity    = 5
	defaultTag         = "logoops"
	expectedDateFormat = "2006-01-02 15:04:05"
	syslogDateFormat   = "2006-01-02T15:04:05.000000-07:00"
)

func getPayload(date, hostname, tag, message string, facility, severity int) string {
	priority := facility*8 + severity
	return fmt.Sprintf("<%d>1 %s %s %s - - %s", priority, date, hostname, tag, message)
}

func main() {
	currentDate := time.Now()
	currentDateAsString := currentDate.Format(expectedDateFormat)

	currentHostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	var host = flag.String("host", defaultHost, "syslog server host")
	var port = flag.Int("port", defaultPort, "syslog server port")
	var tag = flag.String("tag", defaultTag, "syslog tag for the log")
	var date = flag.String("date", currentDateAsString, "date to specify for the message, with format %F %T")
	var hostname = flag.String("hostname", currentHostname, "hostname to associate with the log")
	var message = flag.String("message", defaultMessage, "message to send")
	var facility = flag.Int("facility", defaultFacility, "facility for the log")
	var severity = flag.Int("severity", defaultSeverity, "severity of the log")
	flag.Parse()

	parsedDate, err := time.Parse(expectedDateFormat, *date)
	if err != nil {
		panic(err)
	}
	syslogFormattedDate := parsedDate.Format(syslogDateFormat)

	address := fmt.Sprintf("%s:%d", *host, *port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	payload := getPayload(syslogFormattedDate, *hostname, *tag, *message, *facility, *severity)
	fmt.Fprintf(conn, payload)
	err = conn.Close()
	if err != nil {
		panic(err)
	}
}
