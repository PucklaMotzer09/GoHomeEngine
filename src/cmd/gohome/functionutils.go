package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
)

func isCommandArg(arg string) bool {
	return arg == "build" || arg == "install" || arg == "run" || arg == "generate" || arg == "clean" || arg == "env" || arg == "set" || arg == "reset" || arg == "export"
}

func isValueArg(arg string) bool {
	match, _ := regexp.MatchString("[A-Z]=\\w", arg)
	return match
}

func isConfigArg(arg string) bool {
	return arg == "DEBUG" || arg == "RELEASE"
}

func processValueArg(arg string) (b bool) {
	b = true
	varvalue := strings.Split(arg, "=")
	switch varvalue[0] {
	case "OS":
		VAR_OS = varvalue[1]
	case "ARCH":
		VAR_ARCH = varvalue[1]
	case "FRAME":
		VAR_FRAME = varvalue[1]
	case "RENDER":
		VAR_RENDER = varvalue[1]
	case "START":
		VAR_START = varvalue[1]
	case "CONFIG":
		VAR_CONFIG = varvalue[1]
	case "ANDROID_API":
		VAR_ANDROID_API = varvalue[1]
	case "ANDROID_KEYSTORE":
		VAR_ANDROID_KEYSTORE = varvalue[1]
	case "ANDROID_KEYALIAS":
		VAR_ANDROID_KEYALIAS = varvalue[1]
	case "ANDROID_KEYPWD":
		VAR_ANDROID_KEYPWD = varvalue[1]
	case "ANDROID_STOREPWD":
		VAR_ANDROID_STOREPWD = varvalue[1]
	default:
		CustomValues[varvalue[0]] = varvalue[1]
	}
	return
}

func writeVariables(writer *os.File) {
	if VAR_OS != "" {
		writer.WriteString("OS=" + VAR_OS + "\n")
	}
	if VAR_ARCH != "" {
		writer.WriteString("ARCH=" + VAR_ARCH + "\n")
	}
	if VAR_FRAME != "" {
		writer.WriteString("FRAME=" + VAR_FRAME + "\n")
	}
	if VAR_RENDER != "" {
		writer.WriteString("RENDER=" + VAR_RENDER + "\n")
	}
	if VAR_START != "" {
		writer.WriteString("START=" + VAR_START + "\n")
	}
	if VAR_CONFIG != "" {
		writer.WriteString("CONFIG=" + VAR_CONFIG + "\n")
	}
	if VAR_ANDROID_API != "" {
		writer.WriteString("ANDROID_API=" + VAR_ANDROID_API + "\n")
	}
	if VAR_ANDROID_KEYSTORE != "" {
		writer.WriteString("ANDROID_KEYSTORE=" + VAR_ANDROID_KEYSTORE + "\n")
	}
	if VAR_ANDROID_KEYALIAS != "" {
		writer.WriteString("ANDROID_KEYALIAS=" + VAR_ANDROID_KEYALIAS + "\n")
	}
	if VAR_ANDROID_KEYPWD != "" {
		writer.WriteString("ANDROID_KEYPWD=" + VAR_ANDROID_KEYPWD + "\n")
	}
	if VAR_ANDROID_STOREPWD != "" {
		writer.WriteString("ANDROID_STOREPWD=" + VAR_ANDROID_STOREPWD + "\n")
	}
	for k, v := range CustomValues {
		writer.WriteString(k + "=" + v + "\n")
	}
}

func readVariables(r io.Reader) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("Failed to read", CONFIG_FILE_NAME, ":", err)
		os.Exit(1)
	}
	str := string(contents)
	vars := strings.Split(str, "\n")

	for _, v := range vars {
		varvalue := strings.Split(v, "=")
		switch varvalue[0] {
		case "OS":
			VAR_OS = varvalue[1]
			old_OS = varvalue[1]
		case "ARCH":
			VAR_ARCH = varvalue[1]
			old_ARCH = varvalue[1]
		case "FRAME":
			VAR_FRAME = varvalue[1]
			old_FRAME = varvalue[1]
		case "RENDER":
			VAR_RENDER = varvalue[1]
			old_RENDER = varvalue[1]
		case "START":
			VAR_START = varvalue[1]
			old_START = varvalue[1]
		case "CONFIG":
			VAR_CONFIG = varvalue[1]
			old_CONFIG = varvalue[1]
		case "ANDROID_API":
			VAR_ANDROID_API = varvalue[1]
			old_ANDROID_API = varvalue[1]
		case "ANDROID_KEYSTORE":
			VAR_ANDROID_KEYSTORE = varvalue[1]
			old_ANDROID_KEYSTORE = varvalue[1]
		case "ANDROID_KEYALIAS":
			VAR_ANDROID_KEYALIAS = varvalue[1]
			old_ANDROID_KEYALIAS = varvalue[1]
		case "ANDROID_KEYPWD":
			VAR_ANDROID_KEYPWD = varvalue[1]
			old_ANDROID_KEYPWD = varvalue[1]
		case "ANDROID_STOREPWD":
			VAR_ANDROID_STOREPWD = varvalue[1]
			old_ANDROID_STOREPWD = varvalue[1]
		default:
			if len(varvalue) == 2 {
				CustomValues[varvalue[0]] = varvalue[1]
				old_CustomValues[varvalue[0]] = varvalue[1]
			}
		}
	}
}

func getBuild() Build {
	var varos string
	if VAR_OS == "runtime" {
		varos = runtime.GOOS
	} else {
		varos = VAR_OS
	}

	switch varos {
	case "darwin", "freebsd", "linux", "openbsd", "solaris", "windows":
		return &DesktopBuild{}
	case "android":
		return &AndroidBuild{}
	}

	return nil
}

func resetParameters() {
	// Common
	VAR_OS = "runtime"
	VAR_ARCH = "runtime"
	VAR_FRAME = "GLFW"
	VAR_RENDER = "OpenGL"
	VAR_START = ""
	VAR_CONFIG = "DEBUG"

	// Android
	VAR_ANDROID_API = "16"
	VAR_ANDROID_KEYSTORE = ""
	VAR_ANDROID_KEYALIAS = ""
	VAR_ANDROID_KEYPWD = ""
	VAR_ANDROID_STOREPWD = ""
	CustomValues = make(map[string]string)

	COMMAND = "build"
}

func valuesChanged() bool {
	if VAR_OS != old_OS {
		return true
	}
	if VAR_ARCH != old_ARCH {
		return true
	}
	if VAR_FRAME != old_FRAME {
		return true
	}
	if VAR_RENDER != old_RENDER {
		return true
	}
	if VAR_START != old_START {
		return true
	}
	if VAR_CONFIG != old_CONFIG {
		return true
	}
	if VAR_ANDROID_API != old_ANDROID_API {
		return true
	}
	if VAR_ANDROID_KEYSTORE != old_ANDROID_KEYSTORE {
		return true
	}
	if VAR_ANDROID_KEYALIAS != old_ANDROID_KEYALIAS {
		return true
	}
	if VAR_ANDROID_KEYPWD != old_ANDROID_KEYPWD {
		return true
	}
	if VAR_ANDROID_STOREPWD != old_ANDROID_STOREPWD {
		return true
	}

	for k, v := range CustomValues {
		if v1, ok := old_CustomValues[k]; ok {
			if v != v1 {
				return true
			}
		}
	}

	return false
}