package display

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/achu-1612/glcm"
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
	msg := fmt.Sprintf("\033[31m"+format+"\033[0m", a...)
	fmt.Fprintf(Emitter, "%s", msg)
}

// Successf prints default text to std out.
func Successf(format string, a ...interface{}) {
	msg := fmt.Sprintf("\033[32m"+format+"\033[0m", a...)
	fmt.Fprintf(Emitter, "%s", msg)
}

// Fatalf prints default text to std out and os.Exit(1).
func Fatalf(format string, a ...interface{}) {
	msg := fmt.Sprintf("\033[31m"+format+"\033[0m", a...)
	fmt.Fprintf(Emitter, "%s", msg)
	os.Exit(1)
}

// PrintStatus prints list of service status in tabular format.
func PrintStatus(item *glcm.SocketResponse) {

	out := new(tabwriter.Writer)
	out.Init(Emitter, 0, 8, 1, '\t', 0)

	cols := strings.Split("Name,Status,Uptime,Restarts", ",")
	_, _ = fmt.Fprintln(out, strings.ToUpper(strings.Join(cols, "\t")))

	data := &glcm.RunnerStatus{}

	b, err := json.Marshal(item.Result)
	if err != nil {
		Fatalf("Unable to marshal data, error: %v", err)
	}

	if err := json.Unmarshal(b, data); err != nil {
		Fatalf("Unable to unmarshal data, error: %v", err)
	}

	for name, info := range data.Services {
		var f []string
		f = append(
			f,
			name,
			string(info.Status),
			fmt.Sprintf("%02dh:%02dm:%02ds", int(info.Uptime.Hours()), int(info.Uptime.Minutes())%60, int(info.Uptime.Seconds())%60),
			fmt.Sprintf("%d", info.Restarts),
		)

		_, _ = fmt.Fprintln(out, strings.Join(f, "\t"))
	}
	_ = out.Flush()

	fmt.Println()

	if err != nil {
		Fatalf("Unable to print table, error: %v", err)
	}
}
