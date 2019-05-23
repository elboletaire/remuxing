package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"remuxing/info"
	"strings"
)

func merge() {
	exec.Command(
		"mkvmerge",
		"--title",
		"\"\"",
		"-T",
	)
}

func syntaxError(err string) {
	fmt.Println(fmt.Sprintf("syntax error: %s", err))
	os.Exit(1)
}

func getFileInfo(file string) (info info.Info) {
	fmt.Println("Getting file info for file:", file)

	output, mergeErr := exec.Command(
		"mkvmerge",
		"-F",
		"json",
		"-i",
		file,
	).CombinedOutput()

	if mergeErr != nil {
		fmt.Println(fmt.Sprint(mergeErr) + ": " + string(output))
		return
	}

	if err := json.Unmarshal(output, &info); err != nil {
		panic(err)
	}

	// We want it in minutes, so we don't care about decimals
	info.Container.Properties.Duration = info.Container.Properties.Duration / 1000 / 1000 / 1000

	return
}

func parseArgs() (output string, inputs []string, languages []string) {
	flag.StringVar(&output, "output", "", "The output folder")
	var lang string
	flag.StringVar(&lang, "languages", "", "The desired output languages")

	flag.Parse()

	if len(output) == 0 {
		syntaxError("-output path missing")
	}

	inputs = flag.Args()

	if len(inputs) < 2 {
		syntaxError("at least two inputs are expected")
	}

	languages = strings.Split(lang, ",")

	if len(lang) == 0 || len(languages) < 1 {
		syntaxError("at least one language was expected")
	}

	return
}

func filterInfos(infos []info.Info, test func(info.Info) bool) (ret []info.Info) {
	for _, information := range infos {
		if test(information) {
			ret = append(ret, information)
		}
	}

	return
}

func filterSupported(info info.Info) bool {
	return info.Container.Supported
}

func main() {
	_, inputs, _ := parseArgs()

	var infos []info.Info
	for _, input := range inputs {
		infos = append(infos, getFileInfo(input))
	}
	supported := filterInfos(infos, filterSupported)

	fmt.Printf("%+v\n", supported)

	// fmt.Println(fmt.Sprintf("Output is %s", output))
	// fmt.Println(fmt.Sprintf("Inputs are %s", strings.Join(inputs[:], ",")))
	// fmt.Println(fmt.Sprintf("Languages are %s", strings.Join(languages[:], ",")))
	// args := os.Args[1:]

	// 1. Check expected params are received
	//   - At least two inputs are received
	//   - Output path flag (-o) is also mandatory
}
