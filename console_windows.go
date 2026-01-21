package console

import (
	"os"
	"strings"

	"github.com/UserExistsError/conpty"
	"github.com/runletapp/go-console/interfaces"
)

var _ interfaces.Console = (*consoleWindows)(nil)

type consoleWindows struct {
	initialCols int
	initialRows int

	cpty *conpty.ConPty

	cwd string
	env []string
}

func newNative(cols int, rows int) (Console, error) {
	return &consoleWindows{
		initialCols: cols,
		initialRows: rows,

		cpty: nil,

		cwd: ".",
		env: os.Environ(),
	}, nil
}

func (c *consoleWindows) Start(args []string) error {
	cmdLine := strings.Join(args, " ")

	options := []conpty.ConPtyOption{
		conpty.ConPtyDimensions(c.initialCols, c.initialRows),
		conpty.ConPtyWorkDir(c.cwd),
		conpty.ConPtyEnv(c.env),
	}

	cpty, err := conpty.Start(cmdLine, options...)
	if err != nil {
		return err
	}

	c.cpty = cpty
	return nil
}

func (c *consoleWindows) Read(b []byte) (int, error) {
	if c.cpty == nil {
		return 0, ErrProcessNotStarted
	}

	return c.cpty.Read(b)
}

func (c *consoleWindows) Write(b []byte) (int, error) {
	if c.cpty == nil {
		return 0, ErrProcessNotStarted
	}

	return c.cpty.Write(b)
}

func (c *consoleWindows) Close() error {
	if c.cpty == nil {
		return ErrProcessNotStarted
	}

	c.cpty.Close()
	return nil
}

func (c *consoleWindows) SetSize(cols int, rows int) error {
	c.initialRows = rows
	c.initialCols = cols

	if c.cpty == nil {
		return nil
	}

	return c.cpty.Resize(c.initialCols, c.initialRows)
}

func (c *consoleWindows) GetSize() (int, int, error) {
	return c.initialCols, c.initialRows, nil
}

func (c *consoleWindows) Pid() (int, error) {
	if c.cpty == nil {
		return 0, ErrProcessNotStarted
	}

	return c.cpty.Pid(), nil
}

func (c *consoleWindows) Wait() (*os.ProcessState, error) {
	if c.cpty == nil {
		return nil, ErrProcessNotStarted
	}

	pid := c.cpty.Pid()

	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	return proc.Wait()
}

func (c *consoleWindows) SetCWD(cwd string) error {
	c.cwd = cwd
	return nil
}

func (c *consoleWindows) SetENV(environ []string) error {
	c.env = append(os.Environ(), environ...)
	return nil
}

func (c *consoleWindows) Kill() error {
	if c.cpty == nil {
		return ErrProcessNotStarted
	}

	pid := c.cpty.Pid()

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return proc.Kill()
}

func (c *consoleWindows) Signal(sig os.Signal) error {
	if c.cpty == nil {
		return ErrProcessNotStarted
	}

	pid := c.cpty.Pid()

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return proc.Signal(sig)
}
