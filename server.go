package main

import (
    "fmt"
    "net/http"
    "os"
    "strconv"
    "encoding/json"
)

//IANA port range for dynamic and private connections
const port_min int = 49152
const port_max int = 65535
const einval int = 22

func handler(w http.ResponseWriter, request *http.Request) {
    decoder := json.NewDecoder(request.Body)

    var links []string

    err := decoder.Decode(&links)

    if err != nil {
        fmt.Println("Error on json decoding: ", err)
        w.WriteHeader(400)
        return
    }

    url_info := get_urls_info(links)

    content, err := json.Marshal(url_info)

    if err != nil {
        fmt.Println("Error on json marshaling: ", err)
        w.WriteHeader(500)
        return
    }

    w.Header().Set("Content-type", "application/json")
    w.Write(content)
}

func main() {
    host := os.Getenv(env_host)
    if len(host) == 0 {
        fmt.Printf("Don't know the server host name. " +
                   "Please, set %s enviroment variable\n", env_host)
        os.Exit(-einval)
    }

    str_port := os.Getenv(env_port)
    port, err := strconv.Atoi(str_port)
    if err != nil || port < 49152 || 65636 < port {
        fmt.Printf("Can't get the server port number (got: '%s') " +
                    "Please set %s enviromental variable " +
                    "within range %d and %d\n", str_port, env_port, port_min, port_max)
        os.Exit(-einval)
    }

    bind_to := fmt.Sprintf("%s:%d", host, port)
    fmt.Println("binding server to ", bind_to);

    http.HandleFunc("/", handler)
    err = http.ListenAndServe(bind_to, nil)
    fmt.Println(err)
}
