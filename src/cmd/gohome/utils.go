package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func ExecCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), Env...)
	return cmd.Run()
}

func ConsoleRead() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text[:len(text)-1]
}

func ConsoleReadi() int {
	str := ConsoleRead()
	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0
	} else {
		return int(i)
	}
}

func WorkingDir() string {
	str, _ := os.Getwd()
	return str + "/"
}

func FileExists(fn string) bool {
	_, err := os.Stat(fn)
	return !os.IsNotExist(err)
}

func PackageName() string {
	wd := WorkingDir()
	var slash string
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}
	return wd[strings.LastIndex(wd[:len(wd)-1], slash)+1 : len(wd)-1]
}

func GetSlash() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

func ReplaceStringinFile(fn, old, new string) error {
	file, err := os.Open(fn)
	if err != nil {
		return err
	}

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	file.Close()
	fstr := string(contents)
	fstr = strings.Replace(fstr, old, new, -1)

	newFile, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer newFile.Close()

	if _, err := newFile.WriteString(fstr); err != nil {
		return err
	}

	return nil
}

func GetCustomValue(name string, default_val bool, default_str string) string {
	str, ok := CustomValues[name]
	if !ok {
		if !default_val {
			fmt.Print(name + ": ")
			str = ConsoleRead()
			CustomValues[name] = str
		} else {
			str = default_str
			CustomValues[name] = str
		}
	}

	return str
}

func GetCustomValuei(name string, default_val bool, defaultValue int) int {
	str := GetCustomValue(name, default_val, strconv.Itoa(defaultValue))
	var i1 int
	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		i1 = defaultValue
		CustomValues[name] = strconv.FormatInt(int64(i1), 10)
	} else {
		i1 = int(i)
	}

	return i1
}

func GetCustomValueb(name string, default_val bool, defaultValue bool) bool {
	str := GetCustomValue(name, default_val, strconv.FormatBool(defaultValue))
	var i1 bool
	i, err := strconv.ParseBool(str)
	if err != nil {
		i1 = defaultValue
		CustomValues[name] = strconv.FormatBool(i1)
	} else {
		i1 = i
	}

	return i1
}

func CustomValueExists(name string) (string, bool) {
	v, ok := CustomValues[name]
	return v, ok
}

func AssertValue(val *string, assertstr, text string) {
	if *val == assertstr {
		fmt.Print(text + ": ")
		*val = ConsoleRead()
	}
}

func GetNumberAsWord(number byte) string {
	number -= 48
	switch number {
	case 0:
		return "Zero"
	case 1:
		return "One"
	case 2:
		return "Two"
	case 3:
		return "Three"
	case 4:
		return "Four"
	case 5:
		return "Five"
	case 6:
		return "Six"
	case 7:
		return "Seven"
	case 8:
		return "Eight"
	case 9:
		return "Nine"
	default:
		return "Nan"
	}
}

func LowerCaseAndNoNumber(str string) string {
	str = strings.ToLower(str)
	if str[0] >= 48 && str[0] <= 57 {
		str = strings.Replace(str, string(str[0]), GetNumberAsWord(str[0]), 1)
	}
	str = strings.Replace(str, " ", "_", -1)
	str = strings.Replace(str, "ß", "ss", -1)
	return str
}

func OpenBrowser(url string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	}
}

func FlagSet(flag string) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}

	return false
}
