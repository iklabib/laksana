//go:build linux

/*
originally taken from https://github.com/moby/moby/tree/master/profiles/apparmor

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package apparmor // import "github.com/docker/docker/profiles/apparmor"

import (
	"io"
	"os"
	"path"
	"strings"
	"text/template"
)

// profileDirectory is the file store for apparmor profiles and macros.
const profileDirectory = "/etc/apparmor.d"

// profileData holds information about the given profile for generation.
type profileData struct {
	// Name is profile name.
	Name string
	// DaemonProfile is the profile name of our daemon.
	DaemonProfile string
	// Imports defines the apparmor functions to import, before defining the profile.
	Imports []string
	// InnerImports defines the apparmor functions to import in the profile.
	InnerImports []string
}

// generateDefault creates an apparmor profile from ProfileData.
func (p *profileData) generateDefault(out io.Writer) error {
	compiled, err := template.New("apparmor_profile").Parse(baseTemplate)
	if err != nil {
		return err
	}

	if macroExists("tunables/global") {
		p.Imports = append(p.Imports, "#include <tunables/global>")
	} else {
		p.Imports = append(p.Imports, "@{PROC}=/proc/")
	}

	if macroExists("abstractions/base") {
		p.InnerImports = append(p.InnerImports, "#include <abstractions/base>")
	}

	return compiled.Execute(out, p)
}

// macroExists checks if the passed macro exists.
func macroExists(m string) bool {
	_, err := os.Stat(path.Join(profileDirectory, m))
	return err == nil
}

func GenerateProfile(name string) error {
	p := profileData{
		Name: name,
	}

	// Figure out the daemon profile.
	currentProfile, err := os.ReadFile("/proc/self/attr/current")
	if err != nil {
		// If we couldn't get the daemon profile, assume we are running
		// unconfined which is generally the default.
		currentProfile = nil
	}
	daemonProfile := string(currentProfile)
	// Normally profiles are suffixed by " (enforcing)" or similar. AppArmor
	// profiles cannot contain spaces so this doesn't restrict daemon profile
	// names.
	if parts := strings.SplitN(daemonProfile, " ", 2); len(parts) >= 1 {
		daemonProfile = parts[0]
	}
	if daemonProfile == "" {
		daemonProfile = "unconfined"
	}
	p.DaemonProfile = daemonProfile

	f, err := os.Create(name + ".cfg")
	if err != nil {
		return err
	}

	defer f.Close()

	if err := p.generateDefault(f); err != nil {
		return err
	}

	return nil
}
