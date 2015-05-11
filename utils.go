package main

import (
    "os"
    "log"
    "fmt"
    "path/filepath"
    "strings"
    "regexp"
    "strconv"
)

func ExtractFile(file string, toDir string) []string {
    ext := filepath.Ext(file)
    
    if ext == ".rar" {
        return RARExtractFile(file, toDir)
    } else if ext == ".zip" {
        return ZIPExtractFile(file, toDir)
    }
    
    log.Println("Unknown file type to extract.", filepath.Base(file))
    return []string{}
}

func Exists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func SearchFiles(folders []string) []string {
    fileList := []string{}
    
    for _,path := range folders {
        filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
            base := strings.ToLower(filepath.Base(p))
            ext := filepath.Ext(p)
            file := strings.TrimSuffix(p, ext)
            if !StringInSliceContains(base, config.Ignore_files_with) && StringInSlice(ext, config.Valid_files_ext) && !Exists(file + ".srt") && !Exists(file + ".pt.srt") {
                fileList = append(fileList, p)
            }
            return nil
        })
    }
    
    return fileList
}

func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func StringInSliceContains(a string, list []string) bool {
    for _, b := range list {
        if strings.Contains(a, b) {
            return true
        }
    }
    return false
}

type ParsedFileName struct {
    Filename string
    Title string
    Year int
    Season int
    Episode int
    Group string
    Pack bool
    Movie bool
    TvShow bool
}

func ParseFileName(filename string) ParsedFileName {
    info := ParsedFileName{}
    info.Filename = filename
    
    clean := strings.ToLower(filename)
    for _,i := range config.Clean_from_filename {
        clean = strings.Replace(clean, i, "", -1)
    }
    
    clean = strings.Replace(clean, ".", " ", -1)
    
    rSpaces := regexp.MustCompile(`\s\s+`)
    clean = strings.Trim(rSpaces.ReplaceAllString(clean, " "), " ")
    
    r := regexp.MustCompile(`(.*?)[.\s\-_]([1-2]\d\d\d[.\s\-_])?(us[.\s\-_])?(s(\d+\d)(e(\d+\d)(e\d+\d|\-\d+\d)?)?|(\d+)x(\d+\d))`)
    matches := r.FindStringSubmatch(clean)
    
    if matches != nil {
        info.Title = matches[1]
        info.Year, _ = strconv.Atoi(matches[2])
        info.Season, _ = strconv.Atoi(matches[5])
        info.Episode, _ = strconv.Atoi(matches[7])
        
        if matches[9] != "" && matches[10] != "" {
            info.Season, _ = strconv.Atoi(matches[9])
            info.Episode, _ = strconv.Atoi(matches[10])
        }
        
        if info.Season != 0 && info.Episode == 0 {
            info.Pack = true
        }
        
        info.TvShow = true
    } else {
        r = regexp.MustCompile(`(.*?)([1-2]\d\d\d)`)
        matches = r.FindStringSubmatch(clean)
        if matches != nil && matches[1] != "" {
            info.Title = matches[1]
            info.Year, _ = strconv.Atoi(matches[2])
            info.Movie = true
        }
    }
    
    groups := strings.Split(clean, "-")
    if groups != nil {
        groupsLen := len(groups)
        tmp := strings.Split(strings.Trim(groups[groupsLen-1], " "), " ")
        info.Group = tmp[0]
        clean = strings.Join(groups[0:groupsLen-1], "-")
    }
    
    return info
}

func SearchString(info ParsedFileName) string {
    if info.TvShow {
        return fmt.Sprintf("%s s%02de%02d", info.Title, info.Season, info.Episode)
    } else {
        return info.Title + " " + strconv.Itoa(info.Year)
    }
}
