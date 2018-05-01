package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/femnad/mare"
)

const (
	defaultFacility    = 20
	defaultHost        = "localhost"
	defaultMessage     = "hello"
	defaultPort        = 514
	defaultProtocol    = "tcp"
	defaultSeverity    = 5
	defaultTag         = "logoops"
	expectedDateFormat = "2006-01-02 15:04:05"
	syslogDateFormat   = "2006-01-02T15:04:05.000000-07:00"
)

var validProtocols = []string{"udp", "tcp"}

func getPayload(date, hostname, tag, message string, facility, severity int) string {
	priority := facility*8 + severity
	return fmt.Sprintf("<%d>1 %s %s %s - - %s", priority, date, hostname, tag, message)
}

func ensureCorrectProtocol(protocol string) {
	if !mare.Contains(validProtocols, protocol) {
		errorMessage := fmt.Sprintf("Invalid protocol %s", protocol)
		panic(errorMessage)
	}
}

func parseAndFormatDate(date string) string {
	parsedDate, err := time.Parse(expectedDateFormat, date)
	mare.PanicIfErr(err)
	return parsedDate.Format(syslogDateFormat)
}

func getAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func sendPayload(protocol, address, payload string) {
	conn, err := net.Dial(protocol, address)
	mare.PanicIfErr(err)

	fmt.Fprintf(conn, payload)

	err = conn.Close()
	mare.PanicIfErr(err)
}

func sendMessageViaSyslog(host string, port int, protocol, message, tag, hostname, date string, facility, severity int) {
	ensureCorrectProtocol(protocol)

	syslogFormattedDate := parseAndFormatDate(date)
	address := getAddress(host, port)
	payload := getPayload(syslogFormattedDate, hostname, tag, message, facility, severity)

	sendPayload(protocol, address, payload)
}

func main() {
	currentDate := time.Now()
	currentDateAsString := currentDate.Format(expectedDateFormat)

	currentHostname, err := os.Hostname()
	mare.PanicIfErr(err)

	var date = flag.String("date", currentDateAsString, "date to specify for the message, with format %F %T")
	var facility = flag.Int("facility", defaultFacility, "facility for the log")
	var hostname = flag.String("hostname", currentHostname, "hostname to associate with the log")
	var host = flag.String("host", defaultHost, "syslog server host")
	var message = flag.String("message", defaultMessage, "message to send")
	var port = flag.Int("port", defaultPort, "syslog server port")
	var protocol = flag.String("protocol", defaultProtocol, "protocol for connecting to the server")
	var tag = flag.String("tag", defaultTag, "syslog tag for the log")
	var severity = flag.Int("severity", defaultSeverity, "severity of the log")

	flag.Parse()

	sendMessageViaSyslog(*host, *port, *protocol, *message, *tag, *hostname, *date, *facility, *severity)
}
