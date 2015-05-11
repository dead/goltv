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
                os.MkdirAll(fpath, f.Mode())
            } else {
                var fdir string
                if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
                    fdir = fpath[:lastIndex]
                }
                
                err = os.MkdirAll(fdir, f.Mode())
                if err == nil {
                    fdest, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
                    defer fdest.Close()
                    
                    if err == nil {
                        io.Copy(fdest, rc)
                        ret = append(ret, fpath)
                    }
                }
            }
        }
    }
    
    return ret
}