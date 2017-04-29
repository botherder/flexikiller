// FlexiKiller
// Copyright (C) 2017 Claudio "nex" Guarnieri
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"os"
	"fmt"
	"time"
	"bufio"
	"errors"
	"strings"
	"os/exec"
	"path/filepath"
	"golang.org/x/sys/windows/registry"
	log "github.com/Sirupsen/logrus"
	"github.com/mattn/go-colorable"
)

// These are the Windows services FlexiSpy normally uses.
var flexi_services = []string{
	"ApplicationInitService",
	"ApplicationLookupService",
}

func findFlexiSpy() (string, bool) {
	for _, service := range flexi_services {
		// log.Info("Looking for service with name ", service)

		// Check if the current service exists.
		path := fmt.Sprintf("System\\CurrentControlSet\\Services\\%s", service)
		key_service, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
		if err != nil {
			continue
		}

		// Extract the ImagePath and get the installation folder.
		image_path, _, err := key_service.GetStringValue("ImagePath")
		if strings.Contains(image_path, "Windows Provisioning") {
			return filepath.Dir(image_path), true
		}
	}

	return "", false
}

func disableFlexiSpy() (bool) {
	counter := 0

	for _, service := range flexi_services {
		// Nuke the registry keys that launch the service.
		path := fmt.Sprintf("System\\CurrentControlSet\\Services\\%s", service)
		err := registry.DeleteKey(registry.LOCAL_MACHINE, path)
		if err != nil {
			log.Error(err)
		} else {
			counter += 1
		}
	}

	// If both have been successfully removed, we should be good.
	if counter == 2 {
		return true
	} else {
		return false
	}
}

func uninstallFlexiSpy(base_path string) (error) {
	// Look for the uninstall utility.
	exe_path := filepath.Join(base_path, "uninstall.exe")
	if _, err := os.Stat(exe_path); os.IsNotExist(err) {
		return errors.New("Unable to find uninstall utility")
	}

	// Launching it with argument "clean" seems sufficient to fully
	// uninstall FlexiSpy from the system.
	cmd := exec.Command(exe_path, "clean")
	err := cmd.Start()
	if err != nil {
		log.Error(err)
		return errors.New("Unable to launch uninstall utility")
	}

	// For some reason the FlexiSpy uninstall utility can take a real
	// long time. Just in case, we add a timeout of 15 minutes.
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
		case <-time.After(15 * time.Minute):
			if err := cmd.Process.Kill(); err != nil {
				log.Warning("failed to kill: ", err)
			}
			return errors.New("The uninstall utility is taking too long...")
		case <-done:
			log.Info("Uninstall utility terminated.")
			// We sleep few seconds, just in case.
			time.Sleep(5 * time.Second)
	}

	return nil
}

func finish() {
	log.Info("Press any key to finish ...")
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	os.Exit(0)
}

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())

	log.Println("Looking for records of FlexiSpy...")

	// Check if FlexiSpy is installed.
	flexi_folder, flexi_exists := findFlexiSpy()
	if flexi_exists == false {
		log.Info("I did not find any traces of FlexiSpy.")
		finish()
	}

	log.Warning("Found FlexiSpy at ", flexi_folder)

	// Do we want to uninstall it?
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to uninstall it? (y/N) ")
	choice, _ := reader.ReadString('\n')
	choice = strings.Replace(choice, "\r\n", "", -1)
	if choice == "n" || choice == "" {
		return
	}

	log.Println("Attempting to uninstall FlexiSpy (this might take some time)...")

	// Nuking FlexiSpy.
	err := uninstallFlexiSpy(flexi_folder)
	if err != nil {
		log.Error(err)
	}

	// Check once more if the services are there.
	_, flexi_exists = findFlexiSpy()
	if flexi_exists == false {
		log.Info("FlexiSpy seems to have been removed successfully! :-)")
	// If the damn thing is still there, we remove the Services registry keys.
	// That should be sufficient to disable it at least.
	} else {
		log.Warning("It seems the uninstall somehow failed. :-(")
		log.Println("We try now to disable FlexiSpy, so it is at least not able to survive a reboot.")

		flexi_disabled := disableFlexiSpy()

		if(flexi_disabled) {
			log.Info("FlexiSpy disabled successfully. You might want to restart your computer.")
		} else {
			log.Warning("I did not manage to entirely disable FlexiSpy.")
		}
	}

	finish()
}
