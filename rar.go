package main

import (
    "os/exec"
    "path/filepath"
    "os"
    "log"
)

func RARExtractFile(file string, toDir string) []string {
    log.Println("Running RAR extract", file)

    ret := []string{}
    cmd := exec.Command(config.Unrar, "x", file, toDir)
    err := cmd.Run()
    
    if err != nil {
        log.Println("Could not extract", file, err)
        return ret
    }
    
    
    filepath.Walk(toDir, func(p string, f os.FileInfo, err error) error {
        ret = append(ret, p)
        return nil
    })
    
    return ret
}
