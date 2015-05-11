package main

import (
    "path/filepath"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "bytes"
    "html"
    "regexp"
    "os"
    "io"
    "io/ioutil"
    "strings"
    "log"
)

type SubSearchResult struct {
    id string
    release string
    lang string
}

var jar,_ = cookiejar.New(nil)

func Login() {
    options := cookiejar.Options{}
    jar, _ = cookiejar.New(&options)
    
    client := &http.Client{Jar: jar}
    req, _ := http.NewRequest("POST", "http://legendas.tv/login",  bytes.NewBufferString("_method=POST&data%5BUser%5D%5Busername%5D="+config.Username+"&data%5BUser%5D%5Bpassword%5D="+config.Password+"&data%5Blembrar%5D=on"))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    
    resp, _ := client.Do(req)
    defer resp.Body.Close()
    
    log.Println("Login to LegendasTV.")
}

func SearchSub(search string) []SubSearchResult {
    //defer timeTrack(time.Now(), "Search in LegendasTV for '" + search + "'")
    log.Println("Searching for:", search)
    
    ret := []SubSearchResult{}
    
    client := &http.Client{}
    req, _ := http.NewRequest("GET", "http://legendas.tv/legenda/busca/" + url.QueryEscape(search) + "/1", nil)
    
    req.Header.Set("X-Requested-With", "XMLHttpRequest")
    resp, err := client.Do(req)
    
    if err != nil {
        log.Println("Error: %s\n", err)
    } else {
        r := regexp.MustCompile("<div class=\".*?\">.*?href=\"/download/(.*?)/.*?\".*?>(.*?)</a>.*?<img src=\".*?\" alt=\".*?\" title=\"(.*?)\" /></div>")
        
        defer resp.Body.Close()
        content, err := ioutil.ReadAll(resp.Body)
        
        if err != nil {
            log.Println("Error: %s", err)
        } else {
            matches := r.FindAllStringSubmatch(string(content), -1)
            for i := range matches {
                id := matches[i][1]
                release := html.UnescapeString(matches[i][2])
                lang := matches[i][3]
                
                ret = append(ret, SubSearchResult{id, release, lang})
            }
        }
    }
    
    return ret
}

func DownloadSub(sub SubSearchResult) string {
    //defer timeTrack(time.Now(), "Downloading subtitle "+sub.release)
    
    if Exists(filepath.Join(config.Downloadto, sub.id + ".rar")) {
        return filepath.Join(config.Downloadto, sub.id + ".rar")
    } else if Exists(filepath.Join(config.Downloadto, sub.id + ".zip")) {
        return filepath.Join(config.Downloadto, sub.id + ".zip")
    } else if Exists(filepath.Join(config.Downloadto, sub.id + ".srt")) {
        return filepath.Join(config.Downloadto, sub.id + ".srt")
    } else if Exists(filepath.Join(config.Downloadto, sub.id + ".sub")) {
        return filepath.Join(config.Downloadto, sub.id + ".sub")
    }
    
    client := &http.Client{Jar: jar}
    req, _ := http.NewRequest("GET", "http://legendas.tv/pages/downloadarquivo/" + sub.id, nil)
    resp, err := client.Do(req)
    
    if err != nil {
        log.Println("Error while download sub", sub.id, err)
        return ""
    }
    
    defer resp.Body.Close()
    
    finalUrl := resp.Request.URL.String();
    tokens := strings.Split(finalUrl, "/")
    file := tokens[len(tokens)-1]
    ext := filepath.Ext(file)
    toFile := filepath.Join(config.Downloadto, sub.id + ext)
    
    output, err := os.Create(toFile)
    if err != nil {
        log.Println("Error while creating", toFile, err)
        return ""
    }
    
    _, err = io.Copy(output, resp.Body)
    
    if err != nil {
        log.Println("Error while downloading", toFile, err)
        return ""
    }
    
    return toFile
}