package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {

	//var sepeartor = ","

	var file string
	var separator string
	var textqualifier string

	flag.StringVar(&file, "file", "", "the file")
	flag.StringVar(&separator, "separator", "\r?\n", "the separator")
	flag.StringVar(&textqualifier, "textqualifier", "", "the text qualifier")

	flag.Parse()

	if file == "" {
		log.Fatal("File path not given. Add -file=[file] to flags.")
		return
	}

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("The given file path '", file, "' is invalid")
		} else {
			log.Fatal(err)
		}
	}

	text := string(dat)

	linebreak := regexp.MustCompile("\r?\n")
	separation := regexp.MustCompile(separator)
	sli := make([]string, 0)

	current := ""
	for _, e := range separation.Split(text, -1) {

		e = linebreak.ReplaceAllString(e, "")

		if current != "" {
			current = strings.Join([]string{current, e}, ",")

			if strings.HasSuffix(e, textqualifier) {
				sli = append(sli, strings.Trim(current, textqualifier))
				current = ""
			}
		} else {
			if textqualifier != "" && strings.HasPrefix(e, textqualifier) {
				current = e
				continue
			}

			sli = append(sli, e)
		}

	}
	if current != "" {
		sli = append(sli, current)
	}

	reader := bufio.NewReader(os.Stdin)

	result := make([][]string, 0)

	for _, e := range sli {
		clipboard.WriteAll(e)
		fmt.Print(e, ": ")
		text, _ := reader.ReadString('\n')
		result = append(result, []string{e, linebreak.ReplaceAllString(text, "")})
	}

	fmt.Println(result)

}
