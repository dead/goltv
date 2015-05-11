package main

import (
    "log"
    "path/filepath"
    "strings"
    "os"
    "flag"
)

var configFile string

func main() {
    flag.StringVar(&configFile, "config", "config.json", "path to configuration file")
    flag.Parse()
    
    LoadConfig(configFile)
    Login()
    
    files := SearchFiles(config.Watchdir)
    
    for _, file := range files {
        filename := filepath.Base(file)
        ext := filepath.Ext(file)
        filenameWithoutExt := strings.TrimSuffix(filename, ext)
        
        info := ParseFileName(filenameWithoutExt)
        
        if info.TvShow || info.Movie {
            log.Println("Searching subtitles for", filenameWithoutExt)
            subtitles := SearchSub(SearchString(info))
            subtitle := GetBestSub(subtitles, info)
            
            if subtitle.id != "" {
                log.Println("Found", subtitle)
                
                subFile := DownloadSub(subtitle)
                subFileExt := filepath.Ext(subFile)
                
                if subFileExt == ".rar" || subFileExt == ".zip"  {
                    toSubFile := strings.TrimSuffix(file, ext)
                    extractedFiles := ExtractFile(subFile, filepath.Join(config.Downloadto, subtitle.id) + "/")
                    
                    selectedSubFile := GetBestSubFile(info, extractedFiles)
                    
                    if selectedSubFile != "" {
                        log.Println("Selected subtitle", filepath.Base(selectedSubFile))
                        os.Rename(selectedSubFile, toSubFile + ".pt" + filepath.Ext(selectedSubFile))
                    }
                    
                    os.RemoveAll(filepath.Join(config.Downloadto, subtitle.id) + "/")
                } else if subFileExt == ".srt" || subFileExt == ".sub" {
                    os.Rename(subFile, strings.TrimSuffix(file, ext)+subFileExt)
                } else {
                    log.Println("Unknown subtitle extension", subFileExt)
                }
            } else {
                log.Println("Subtitle not found")
            }
        }
    }
}

func GetBestSub(subtitles []SubSearchResult, toFileInfo ParsedFileName) SubSearchResult {
    ret := SubSearchResult{}
    
    if toFileInfo.TvShow {
        for _, subtitle := range subtitles {
            subtitleInfo := ParseFileName(subtitle.release)
            release := strings.ToLower(subtitle.release)
            if toFileInfo.Title == subtitleInfo.Title || strings.Contains(subtitleInfo.Title, toFileInfo.Title) || strings.Contains(toFileInfo.Title, subtitleInfo.Title) {
                if subtitleInfo.Year == 0 || toFileInfo.Year == 0 || toFileInfo.Year == subtitleInfo.Year {
                    if subtitleInfo.Season == toFileInfo.Season {
                        if subtitleInfo.Pack || toFileInfo.Episode == subtitleInfo.Episode {
                            if toFileInfo.Group == "" || strings.Contains(release, toFileInfo.Group) {
                                return subtitle
                            } else {
                                ret = subtitle
                            }
                        }
                    }
                }
            }
        }
    } else if toFileInfo.Movie {
        for _, subtitle := range subtitles {
            subtitleInfo := ParseFileName(subtitle.release)
            release := strings.ToLower(subtitle.release)
            
            if toFileInfo.Title == subtitleInfo.Title || strings.Contains(subtitleInfo.Title, toFileInfo.Title) || strings.Contains(toFileInfo.Title, subtitleInfo.Title) {
                if toFileInfo.Year == subtitleInfo.Year {
                    if toFileInfo.Group == "" || strings.Contains(release, toFileInfo.Group) {
                        return subtitle
                    } else {
                        if strings.Contains(toFileInfo.Filename, "720p") && strings.Contains(release, "720p") {
                            ret = subtitle
                        } else if strings.Contains(toFileInfo.Filename, "1080p") && strings.Contains(release, "1080p") {
                            ret = subtitle
                        } else if !strings.Contains(toFileInfo.Filename, "720p") && !strings.Contains(toFileInfo.Filename, "1080p") {
                            ret = subtitle
                        }
                    }
                }
            }
        }
    }
    
    return ret
}

func GetBestSubFile(info ParsedFileName, files []string) string {
    ret := ""
    
    for _, file := range files {
        ext := filepath.Ext(file)
        name := filepath.Base(file)
        nameWithoutExt := strings.TrimSuffix(name, ext)
        fInfo := ParseFileName(nameWithoutExt)
        
        if fInfo.TvShow == info.TvShow && fInfo.Movie == info.Movie {
            if fInfo.TvShow && fInfo.Episode == info.Episode && fInfo.Season == info.Season {
                if strings.Contains(strings.ToLower(name), info.Group) {
                    return file
                } else {
                    ret = file
                }
            } else if fInfo.Movie && fInfo.Year == info.Year {
                if strings.Contains(strings.ToLower(name), info.Group) {
                    return file
                } else {
                    ret = file
                }
            }
        }
    }
    
    return ret
}