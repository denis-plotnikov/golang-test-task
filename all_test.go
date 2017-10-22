package main

import (
    "testing"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "bytes"
    "io/ioutil"
)


func send_links(query []byte, t *testing.T) (int, []byte, error) {

    host := os.Getenv(env_host)

    port, err := strconv.Atoi(os.Getenv(env_port))

    if err != nil {
        t.Errorf("Error on port aquiring: ", err.Error())
        return 0, nil, err
    }

    send_to := fmt.Sprintf("http://%s:%d", host, port)

    req, err := http.NewRequest("POST", send_to, bytes.NewBuffer(query))

    if err != nil {
        t.Errorf("Error on http req creating: ", err.Error())
        return 0, nil, err
    }

    req.Header.Set("Content-Type", "application/json")

    sender := &http.Client{}

    res, err := sender.Do(req)

    if err != nil {
        t.Errorf("Error on request sending: ", err.Error())
        return 0, nil, err
    }
    defer res.Body.Close()

    var body []byte

    if res.StatusCode == 200 {
        // body, err := ioutil.ReadAll(res.Body)
        _, err := ioutil.ReadAll(res.Body)
        if err != nil {
            return res.StatusCode, nil, err
        }
        // fmt.Println("response Body:")
        // fmt.Println(string(body))
    }

    return res.StatusCode, body, nil

}

func TestSendValid(t *testing.T) {
    var link_list []string;
    link_list = append(link_list, "http://example.com")
    link_list = append(link_list, "https://example.com")

    query, err := json.Marshal(link_list)

    if err != nil {
        t.Errorf("Error on json encoding: ", err.Error())
        return
    }

    retcode, _, err := send_links(query, t)

    if err != nil {
        return
    }

    if retcode != 200 {
        t.Errorf("Error on response status. Expected: 200 Got:", retcode)
    }
}

func TestSendNotValid(t *testing.T) {
    var link_list []string;
    link_list = append(link_list, "http://the-middle-of-nowhere.com")
    link_list = append(link_list, "https://horns-and-hooves.biz")

    query, err := json.Marshal(link_list)
    if err != nil {
        t.Errorf("Error on json encoding: ", err.Error())
        return
    }

    retcode, _, err := send_links(query, t)

    if err != nil {
        return
    }

    if retcode != 200 {
        t.Errorf("Error on response status. Expected: 200 Got:", retcode)
    }
}

func TestSendWrongJson(t *testing.T) {
    retcode, _, err := send_links([]byte{0}, t)
    if err != nil {
        return
    }

    if retcode != 400 {
        t.Errorf(
            "Error on response status. Expected: 400 Got: ", retcode)
    }
}
