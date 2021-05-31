package systemd

import (
	"os/exec"
)

// daemonReload performs a 'systemctl daemon-reload'
func daemonReload() ([]byte, error) {
    cmd := exec.Command("systemctl", "daemon-reload")
    return cmd.CombinedOutput()
}

// disableService disables the specified systemd service unit
func disableService(svc string) ([]byte, error) {
	cmd := exec.Command("systemctl", "disable", svc)
	return cmd.CombinedOutput()
}

// enableService enables the specified systemd service unit
func enableService(svc string) ([]byte, error) {
	cmd := exec.Command("systemctl", "enable", svc)
	return cmd.CombinedOutput()
}

// reloadService reloads the specified systemd service unit
// NB: Not all services support this command.  It is common in web servers.
func reloadService(svc string) ([]byte, error) {
    cmd := exec.Command("systemctl", "reload", svc)
    return cmd.CombinedOutput()
}

// startService starts the specified systemd service unit
func startService(svc string) ([]byte, error) {
    cmd := exec.Command("systemctl", "start", svc)
    return cmd.CombinedOutput()
}

// stopService stops the specified systemd service unit
func stopService(svc string) ([]byte, error) {
    cmd := exec.Command("systemctl", "stop", svc)
    return cmd.CombinedOutput()
}