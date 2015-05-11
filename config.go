package main

import (
    "encoding/json"
    "os"
    "log"
)

type Configuration struct {
    Unrar                   string
    Downloadto              string
    Watchdir                []string
    Username                string
    Password                string
    Valid_files_ext         []string
    Ignore_files_with       []string
    Clean_from_filename     []string
}

var config = Configuration{}

func LoadConfig(file string) {
    f, _ := os.Open(file)
    d := json.NewDecoder(f)
    err := d.Decode(&config)
    if err != nil {
        log.Fatal(err)
    }
}