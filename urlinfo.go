package main

import (
    "golang.org/x/net/html"
    "io/ioutil"
    "net/http"
    "bytes"
    "sort"
    "time"
)

const request_timeout int = 1000; // in milisec

type TagCount map[string] int

type UrlMetadata struct {
    Status int            `json:"status"`
    Content_type *string  `json:"content-type,omitempty"`
    Content_length *int   `json:"content-length,omitempty"`
}

type TagInfo struct {
    Tag_name string       `json:"tag-name"`
    Count int             `json:"count"`
}

type UrlInfo struct {
    Url string            `json:"url"`
    Meta UrlMetadata      `json:"meta"`
    Elements *[]TagInfo   `json:"elements,omitempty"`
}

func get_url_data(url string) (*UrlMetadata, []byte) {
    var meta UrlMetadata

    timeout := time.Duration(time.Duration(request_timeout) * time.Millisecond)
    client := http.Client{ Timeout: timeout }
    resp, err := client.Get(url)

    if err != nil {
        return &meta, nil
    }

    status := resp.StatusCode

    meta.Status = status

    if status < 200 || 226 < status {
        return &meta, nil
    }

    defer resp.Body.Close()
    var content_type = resp.Header.Get("Content-type")
    meta.Content_type = &content_type

    content, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return &meta, nil
    }

    l := len(content)
    meta.Content_length = &l

    return &meta, content
}

func count_tags(content []byte) TagCount {
    // TODO: if content is too big (threshold?) we could split it
    //       into equal parts, process the parts in goroutines and
    //       then merge the results into the same map.
    //       Do it later in case of latency problems
    m := make(TagCount)
    tokenizer := html.NewTokenizer(bytes.NewReader(content))

    for {
        token := tokenizer.Next()
        switch {
            case token == html.StartTagToken || token == html.SelfClosingTagToken:
                t := tokenizer.Token()
                tag_name := t.Data
                _, ok := m[tag_name]
                if ok {
                    m[tag_name] += 1
                } else {
                    m[tag_name] = 1
                }
            case token == html.ErrorToken:
                return m
        }
    }
}

func get_sorted_elements(tc TagCount) *[]TagInfo {
   var tags []TagInfo

   // sort the elements in alfabetical order
   // sorting can add some mem and cpu overhead but not much
   // shouldn't be a problem -- redo if optimization needed
   // Actually, it's here just in case somebody wants to watch the results
   // with his/her eyes. Leave it here just to show off a little bit
   var keys []string

   for k, _ := range tc {
       keys = append(keys, k)
   }

   sort.Strings(keys)

   for _, key := range keys {
       tags = append(tags, TagInfo{ Tag_name: key, Count: tc[key] })
   }

   return &tags
}

func get_url_info(url string, ch_res chan<-UrlInfo) {
    meta, content := get_url_data(url)

    info := UrlInfo{ Url: url, Meta: *meta }

    if content != nil {
        tag_count := count_tags(content)
        info.Elements = get_sorted_elements(tag_count)
    }

    ch_res <- info
}


func get_urls_info(urls []string) []UrlInfo {
    var urls_info []UrlInfo

    ch_res := make(chan UrlInfo)

    for _, url := range urls {
        go get_url_info(url, ch_res)
    }

    urls_number := len(urls)
    done_counter := 0

    for {
        select {
        case info := <-ch_res:
            urls_info = append(urls_info, info)
            done_counter++
            if done_counter == urls_number {
                return urls_info
            }
       }
   }
}
