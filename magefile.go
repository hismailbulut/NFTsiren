//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

const (
	OS_LINUX   = "linux"
	OS_WINDOWS = "windows"
	OS_DARWIN  = "darwin"
	OS_ANDROID = "android"
	OS_IOS     = "ios"
)

const (
	EXE = ".exe" // for windows
	APP = ".app" // for darwin
	APK = ".apk" // for android
	AAB = ".aab" // for android distribution
)

const (
	ENV_GOOS             = "GOOS"
	ENV_JAVA_HOME        = "JAVA_HOME"
	ENV_JAVA8_HOME       = "JAVA8_HOME"
	ENV_ANDROID_HOME     = "ANDROID_HOME"
	ENV_ANDROID_NDK_HOME = "ANDROID_NDK_HOME"
	ENV_ANDROID_SDK_ROOT = "ANDROID_SDK_ROOT"
)

// Utility functions

func sizeoffile(path string) int64 {
	f, err := os.Stat(path)
	if err != nil {
		panic(fmt.Errorf("os.Stat failed: %s", err))
	}
	if f.IsDir() {
		panic(fmt.Errorf("couldn't calculate size of directory"))
	}
	return f.Size()
}

func getarg(index int) string {
	if index >= 0 && index+2 < len(os.Args) {
		return os.Args[index+2]
	}
	return ""
}

func targetfromargs(index int) string {
	target := getarg(0)
	if target == "" {
		target = runtime.GOOS
	}
	return target
}

func cwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("couldn't find cwd: %s", err))
	}
	return wd
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Errorf("file exists but there is error: %s", err))
		}
		return false
	}
	return true
}

func join(elem ...string) string {
	return filepath.Join(elem...)
}

func joinargs(elem ...string) string {
	return strings.Join(elem, " ")
}

func abs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("failed to find absolute path of %s because: %s", path, err))
	}
	return abs
}

func run(env map[string]string, cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = os.Environ()
	for k, v := range env {
		// search and replace or append
		index := -1
		for i, e := range c.Env {
			if k == strings.Split(e, "=")[0] {
				index = i
			}
		}
		envstr := fmt.Sprintf("%s=%s", k, v)
		if index < 0 {
			c.Env = append(c.Env, envstr)
		} else {
			c.Env[index] = envstr
		}
	}
	fmt.Println("RUN:", cmd, strings.Join(args, " "))
	err := c.Run()
	if err != nil {
		panic(fmt.Errorf("failed to run command '%s': %s", cmd, err))
	}
}

func mkdir(path string) {
	fmt.Println("MKDIR:", path)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		panic(fmt.Errorf("mkdir failed: %s", err))
	}
}

func chdir(dir string) {
	fmt.Println("CHDIR:", dir)
	err := os.Chdir(dir)
	if err != nil {
		panic(fmt.Errorf("chdir failed: %s", err))
	}
}

func move(src, dst string) {
	fmt.Println("MOVE:", src, dst)
	err := os.Rename(src, dst)
	if err != nil {
		panic(fmt.Errorf("failed to move file %s to %s: %s", src, dst, err))
	}
}

func ensureGoCompilerInstalled() {
	_, err := exec.LookPath("go")
	if err != nil {
		panic(fmt.Errorf("go compiler is not installed: %s", err))
	}
	// check the compiler version
	run(nil, "go", "version")
}

// Repo must end with version tag or @latest
func ensureGoAppInstalled(name, repo string) {
	_, err := exec.LookPath(name)
	if err != nil {
		// TODO: go get at root repo
		run(nil, "go", "install", repo)
	}
}

// NFTSIREN APP DEFINES
const (
	NAME        = "nftsiren"
	SOURCE_DIR  = "./cmd/nftsiren"
	VERSION     = "0.0.1"
	BIN_DIR     = "bin"
	DEBUG_DIR   = "bin/debug"
	RELEASE_DIR = "bin/release"
)

func ensureOutputPaths() {
	if !exists(DEBUG_DIR) {
		mkdir(DEBUG_DIR)
	}
	if !exists(RELEASE_DIR) {
		mkdir(RELEASE_DIR)
	}
}

func getVersionNumber() int {
	onlynum := ""
	for _, c := range VERSION {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			onlynum += string(c)
		}
	}
	num, err := strconv.Atoi(onlynum)
	if err != nil {
		panic(err)
	}
	return num
}

func build(release bool) {
	ensureGoCompilerInstalled()
	ensureOutputPaths()
	if release {
		ensureGoAppInstalled("gogio", "gioui.org/cmd/gogio@latest")
	}
	target := targetfromargs(0)
	output := NAME
	env := map[string]string{}
	needsPatching := false
	switch target {
	case OS_WINDOWS:
		output = NAME + EXE
		if release {
			needsPatching = true
		}
	case OS_DARWIN:
		if release {
			target = "macos"
			output = NAME + APP
			needsPatching = true
		}
	case OS_LINUX:
		// We don't patch on linux
	case OS_ANDROID:
		fmt.Println(ENV_ANDROID_HOME+":", os.Getenv(ENV_ANDROID_HOME))
		fmt.Println(ENV_ANDROID_NDK_HOME+":", os.Getenv(ENV_ANDROID_NDK_HOME))
		fmt.Println(ENV_JAVA_HOME+":", os.Getenv(ENV_JAVA8_HOME))
		env[ENV_ANDROID_SDK_ROOT] = os.Getenv(ENV_ANDROID_HOME)
		env[ENV_JAVA_HOME] = os.Getenv(ENV_JAVA8_HOME)
		if release {
			output = NAME + AAB
		} else {
			output = NAME + APK
		}
		// We always patch on android
		needsPatching = true
	case OS_IOS:
		// TODO
		panic("ios builds not implemented")
	}

	var outdir string
	ldflags := []string{
		fmt.Sprintf("-X main.VERSION=%s", VERSION),
	}
	tags := []string{}
	if release {
		ldflags = append(ldflags, "-w", "-s")
		outdir = RELEASE_DIR
		tags = append(tags, "release")
	} else {
		outdir = DEBUG_DIR
		tags = append(tags, "debug")
	}

	if needsPatching {
		run(env, "gogio",
			"-target", target,
			"-o", join(outdir, output),
			"-icon", join(SOURCE_DIR, "assets/images/icon.png"),
			"-appid", "com.nftsiren.app",
			"-version", strconv.Itoa(getVersionNumber()),
			"-ldflags", joinargs(ldflags...),
			"-tags", joinargs(tags...),
			"-x", SOURCE_DIR,
		)
	} else {
		env[ENV_GOOS] = target
		args := []string{"build",
			"-o", join(outdir, output),
			"-tags", joinargs(tags...),
			"-ldflags", joinargs(ldflags...),
		}
		if !release {
			args = append(args, "-race")
		}
		args = append(args, SOURCE_DIR)
		run(env, "go", args...)
	}
}

func Build() {
	build(false)
}

func Run() {
	mg.Deps(Build)
	target := targetfromargs(0)
	output := NAME
	switch target {
	case OS_WINDOWS:
		output = NAME + EXE
	}
	run(nil, join(DEBUG_DIR, output))
}

func Release() {
	build(true)
}

func Tests() {
	ensureGoCompilerInstalled()
	run(nil, "go", "test", "./...")
}

func Staticcheck() {
	ensureGoAppInstalled("staticcheck", "honnef.co/go/tools/cmd/staticcheck@latest")
	run(nil, "staticcheck", "./...")
}
