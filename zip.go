package main

import (
    "os"
    "io"
    "archive/zip"
    "path/filepath"
    "log"
    "strings"
)

func ZIPExtractFile(file string, toDir string) []string {
    log.Println("Running ZIP extract", file)

    ret := []string{}
    r, err := zip.OpenReader(file)
    
    if err != nil {
        log.Println("Could not open", file, err)
        return ret
    }
    
    defer r.Close()
    
    for _, f := range r.File {
        rc, err := f.Open()
        
        if err == nil {
            fpath := filepath.Join(toDir, f.Name)
            
            if f.FileInfo().IsDir() {
                os.MkdirAll(fpath, 0777)
            } else {
                var fdir string
                if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
                    fdir = fpath[:lastIndex]
                }
                
                err = os.MkdirAll(fdir, 0777)
                if err == nil {
                    fdest, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
                    defer fdest.Close()
                    
                    if err == nil {
                        log.Println("Extract: ", fpath)
                        io.Copy(fdest, rc)
                        ret = append(ret, fpath)
                    } else {
                        log.Println("Could not extract", fpath, err)
                    }
                } else {
                    log.Println("Could create dir to extract", fdir, err)
                }
            }
        }
    }
    
    return ret
}