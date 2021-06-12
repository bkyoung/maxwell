package systemd

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "strings"
)

// unit represents a systemd unit file
type unit struct {
    name           string
    unitPath       string
    content        []byte
    disabled       bool
    executablePath string
}

// opt is a functional configuration option for configuring a systemd unit
type opt func(u *unit)

// New returns a unit, which can be operated on like a systemd unit
func New(name string, opts ...opt) *unit {
    u := &unit{
        name: name,
    }

    for _, o := range opts {
        o(u)
    }

    return u
}

// Configure performs additional configuration of a systemd unit
func (u *unit) Configure(opts ...opt) {
    for _, o := range opts {
        o(u)
    }
}

// WithExecutablePath allows for (re)configuration of a unit's ExecStart path
func WithExecutablePath(path string) opt {
    return func(u *unit) {
        u.executablePath = path
    }
}


// WithUnitPath allows for (re)configuration of a unit's filesystem location on disk
func WithUnitPath(path string) opt {
    return func(u *unit) {
        u.unitPath = path
    }
}

// WithUnitContent allows for creating or changing the content of a unit file
func WithUnitContent(c []byte) opt {
    return func(u *unit) {
        u.content = c
    }
}

// WithUnitDisabled allows for disabling a unit, which one may want to do during installation
// This allows for overriding the default behavior in Ubuntu, if desired, which is to
// enable and start a unit during installation.
func WithUnitDisabled(d bool) opt {
    return func(u *unit) {
        u.disabled = d
    }
}

// Install installs a systemd unit to disk, performs 'systemctl daemon-reload' as needed
func Install(u *unit) error {
    written, err := writeFile(u.content, u.unitPath);if err != nil {
        return fmt.Errorf("Error writing unit file: %w\n", err)
    }
    if written {
        out, err := daemonReload();if err != nil {
            return fmt.Errorf("Error performing 'systemctl daemon-reload': \n%s\n%w\n", string(out), err)
        }
    }
    return nil
}

// Uninstall removes the installed systemd unit from disk
func Uninstall(u *unit) error {
    if checkExists(u.unitPath) {
        deleted, err := deleteFile(u.unitPath)
        if err != nil {
            return fmt.Errorf("Error deleting unit file: %w\n", err)
        }
        if deleted {
            out, err := daemonReload()
            if err != nil {
                return fmt.Errorf("Error performing 'systemctl daemon-reload': \n%s\n%w\n", string(out), err)
            }
        } else {
            return fmt.Errorf("%s was not deleted", u.unitPath)
        }
    }
    return nil
}

// checkExists returns true if absolute path f exists on disk
func checkExists(f string) bool {
    _, err := os.Stat(f);if err == nil {
        return true
    }
    return false
}

// checkMatches returns true if the content of a and b are identical
func checkMatches(a, b []byte) bool {
    res := bytes.Compare(a, b)
    if res == 0 {
        return true
    }
    return false
}

// deleteFile deletes the file located at the specified path on disk
// returns true if unit file existed and is successfully deleted
func deleteFile(path string) (bool, error) {
    if checkExists(path) {
        err := os.Remove(path);if err != nil {
            return false, fmt.Errorf("Error deleting unit fie: %w\n", err)
        }
        return true, nil
    }
    return false, nil
}

// readFile returns the content of the specified file located at path
func readFile(path string) ([]byte, error) {
    return os.ReadFile(path)
}

// execPathFromUnit extracts the executable's filesystem path from an existing service unit
func execPathFromUnit(u []byte) (string, error) {
    var execPath string
    b := bytes.NewBuffer(u)
    scanner := bufio.NewScanner(b)
    scanner.Split(bufio.ScanLines)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "ExecStart") {
            execPath = strings.Split(line, "=")[1]
        }
    }
    return execPath, nil
}

// writeFile writes the supplied content to the specified path on disk
// returns true only if bytes are written to disk (i.e. file changes)
func writeFile(file []byte, path string) (bool, error) {
    if checkExists(path) {
        content, err := readFile(path);if err != nil {
            return false, fmt.Errorf("Error reading unit file: %w\n", err)
        }
        if checkMatches(file, content) {
            return false, nil
        }
    }

    err := os.WriteFile(path, file, 0755)
    return true, err
}
