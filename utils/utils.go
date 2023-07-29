package utils

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type Lesson struct {
	HomeWorks map[string]string
}

var HomeworkMap = map[int]Lesson{}

func UpdateCache() {
	go func() {
		b := HomeworkMap
		print(b)
		data, _ := json.MarshalIndent(HomeworkMap, "", " ")
		file, err := os.OpenFile("./files/data.json", os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("Failed to open file: %v", err)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Printf("Failed to write to file: %v", err)
			}
		}(file)

		_, err = file.Write(data)
		if err != nil {
			log.Printf("Failed to write to file: %v", err)
		}
	}()
}

func PopulateFromCache() {
	go func() {
		file, err := os.OpenFile("./files/data.json", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&HomeworkMap)
		if err != nil {
			if errors.Is(err, io.EOF) {
				HomeworkMap = map[int]Lesson{}
			}
		}
	}()
}
