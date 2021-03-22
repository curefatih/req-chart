package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
>>	Rar içeriğinin bulunduğu dizin (relative)
*/
var DIR = "./resources"

func renameFile(dirWithName, new string) {
	os.Rename(dirWithName, new)
}

func main() {
	files, err := ioutil.ReadDir(DIR)
	if err != nil {
		log.Fatal(err)
	}

	var total = make(map[string]string)

	for _, file := range files {

		da, _ := base64.StdEncoding.DecodeString(file.Name())

		content, err := ioutil.ReadFile(DIR + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		res := strings.Split(string(content), " ")

		word := ""
		for _, v := range res {
			char, _ := strconv.ParseInt(v, 2, 64)
			ch := fmt.Sprintf("%c", char)
			word += ch
		}

		total[string(da)] = word
	}

	for i := 0; i < len(files); i++ {
		str := fmt.Sprintf("%d", i)
		fmt.Printf("%s", total[str])
	}

}
