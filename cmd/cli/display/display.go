package display

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/achu-1612/glcm"
	"github.com/deis/deis/pkg/prettyprint"
)

var Emitter io.Writer = os.Stdout

// Printf prints the response from the socket.
func Printf(r *glcm.SocketResponse) {
	if r.Status == glcm.Success {
		Successf("%v\n", r.Result)
	} else {
		Errorf("%v\n", r.Result)
	}
}

// Errorf prints default text to std out.
func Errorf(format string, a ...interface{}) {
	format = prettyprint.Colorize(fmt.Sprintf("{{.Red}}%s{{.Default}}", format))
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(Emitter, "%s", msg)
}

// Successf prints default text to std out.
func Successf(format string, a ...interface{}) {
	format = prettyprint.Colorize(fmt.Sprintf("{{.Green}}%s{{.Default}}", format))
	fmt.Fprintf(Emitter, format, a...)
}

// Fatalf prints default text to std out and os.Exit(1).
func Fatalf(format string, a ...interface{}) {
	format = prettyprint.Colorize(fmt.Sprintf("{{.Red}}%s{{.Default}}", format))
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(Emitter, "%s", msg)
	os.Exit(1)
}

// PrintStatus prints list of service status in tabular format.
func PrintStatus(item *glcm.SocketResponse) {

	out := new(tabwriter.Writer)
	out.Init(Emitter, 0, 8, 1, '\t', 0)

	data := &glcm.RunnerStatus{}

	b, err := json.Marshal(item.Result)
	if err != nil {
		Fatalf("Unable to marshal data, error: %v", err)
	}

	if err := json.Unmarshal(b, data); err != nil {
		Fatalf("Unable to unmarshal data, error: %v", err)
	}

	for k, v := range data.Services {
		var f []string
		f = append(f, prettyprint.Colorize(fmt.Sprintf("{{.Yellow}}%s{{.Default}}", k)))

		switch v.Status {
		case glcm.ServiceStatusRegistered, glcm.ServiceStatusScheduled, glcm.ServiceStatusScheduledForRestart:
			f = append(f, prettyprint.Colorize(fmt.Sprintf("{{.Blue}}%s{{.Default}}", v.Status)), v.Uptime.String())

		case glcm.ServiceStatusRunning:
			f = append(f, prettyprint.Colorize(fmt.Sprintf("{{.Green}}%s{{.Default}}", v.Status)), v.Uptime.String())

		case
			glcm.ServiceStatusStopped, glcm.ServiceStatusExhausted, glcm.ServiceStatusExited:
			f = append(f, prettyprint.Colorize(fmt.Sprintf("{{.Red}}%s{{.Default}}", v.Status)), v.Uptime.String())
		}

		_, _ = fmt.Fprintln(out, strings.Join(f, "\t"))
	}
	_ = out.Flush()

	fmt.Println()

	if err != nil {
		Fatalf("Unable to print table, error: %v", err)
	}
}
