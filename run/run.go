// Package run has two main functions:
// - work as a task runner, watching files and
// rebuilding them if necessary;
// - works as a proxy server, that runs a user application,
// proxies all requests to it, and shows detailed
// error messages if needed (TODO).
package run

import (
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/anonx/sunplate/command"
	"github.com/anonx/sunplate/log"
	p "github.com/anonx/sunplate/path"
	"github.com/anonx/sunplate/watcher"

	"github.com/tsuru/config"
)

// ConfigFile is a name of the file that is located at the
// root of user project and describes what the test runner should do.
var ConfigFile = "sunplate.yml"

// Handler is an instance of run subcommand.
var Handler = command.Handler{
	Name:  "run",
	Info:  "start a task runner",
	Usage: "run [path]",
	Desc: `
`,

	Main: start,
}

var (
	// started indicates whether user app is started.
	started bool

	// stopChannel and stopped are channels that are used to let us know
	// the user app has to be stopped or has already been stopped,
	// so we can safely start it again.
	stopChannel = make(chan bool, 1)
	stopped     = make(chan bool, 1)
)

// start is an entry point of the command.
var start = func(action string, params command.Data) {
	imp := p.AbsoluteImport(params.Default(action, "./"))
	dir := p.PackageDir(imp)
	cf := filepath.Join(dir, ConfigFile)

	// Execute all commands from the requested directory.
	os.Chdir(dir)

	// Trying to read a configuration file..
	err := config.ReadConfigFile(cf)
	log.AssertNil(err)

	// Execute all init tasks if they exist.
	is, err := config.GetList("init")
	if err == nil {
		execute(is)
	}

	// Get after tasks if they exist.
	after, _ := config.GetList("after")

	// Build and start the user app for the first time.
	execute(after)
	run(userCommand(imp))

	// Extract patterns and tasks from watch section of config file.
	// And add them to watcher.
	w := watcher.NewType()
	ws, err := config.Get("watch")
	log.AssertNil(err)
	m := ws.(map[interface{}]interface{})
	for k := range m {
		p := k.(string)

		ts, err := config.GetList("watch:" + p)
		log.AssertNil(err)

		w.Listen(p, func() {
			execute(ts)
			execute(after)
			run(userCommand(imp))
		})
	}

	// Cleaning up after we are done.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Warn.Fatal("Application has been stopped.")
}

// userCommand returns a command that should be used for
// starting user application.
func userCommand(imp string) string {
	return filepath.Base(imp)
}

// execute gets a list of tasks and starts them.
func execute(tasks []string) {
	// Iterate over all available tasks.
	for i := range tasks {
		n, as := task(tasks[i])
		cmd := exec.Command(n, as...)
		bs, err := cmd.Output()
		if err != nil {
			log.Error.Panicf(`Failed to execute a command "%s", error: %v.`, tasks[i], err)
		}
		log.Info.Printf("`%s`:\n%s", tasks[i], bs)
	}
}

// run starts a new instance of a task. At the same time
// its previous instance is stopped.
func run(t string) {
	// Stopping previous instance of the app.
	if started {
		stopChannel <- true
		<-stopped
	}

	// Show message about starting a new instance.
	started = true

	// Parse the input task, prepare a command.
	n, as := task(t)
	cmd := exec.Command(n, as...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Starting the user task.
	err := cmd.Start()
	if err != nil {
		log.Error.Panicf("Failed to start a command `%s`, error: %v.", t, err)
	}

	// Make sure we'll be able to stop the app.
	go func() {
		<-stopChannel

		cmd.Process.Kill()
		_, err := cmd.Process.Wait()
		log.AssertNil(err)

		stopped <- true
		started = false
	}()
}

// task gets a string representation and returns
// a name of the command and arguments.
func task(s string) (string, []string) {
	ps := strings.Split(s, " ") // No spaces are allowed between arguments.
	// We are not checking the length of ps as
	// a garanteed minimum is 1.
	// tsuru/config returns <nil> instead of empty string.
	var as []string
	if len(ps) > 1 { // If the task has not only command but also arguments.
		as = ps[1:]
	}
	return ps[0], as
}