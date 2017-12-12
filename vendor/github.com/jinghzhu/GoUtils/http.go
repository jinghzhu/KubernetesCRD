package utils


import (
    "bytes"
    "fmt"
    "errors"
    "io"
    "io/ioutil"
    "regexp"
    "strings"
    "net/http"
)


const (
    SUCCESS_HTTP_CODE  = 0
    HTTP_METHOD_GET    = "GET"
    HTTP_METHOD_POST   = "POST"
    HTTP_METHOD_DELETE = "DELETE"
    HTTP_METHOD_PUT    = "PUT"
)


func IsEmail(email string) bool {
    email = strings.TrimSpace(email)
    errMsgRegeXp := "Error in running regular expression"
    urlRegeXp := "^([a-zA-Z0-9_\\-\\.]+)@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3\\}\\.)|(([a-zA-Z0-9\\-]+\\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\\]?)$"
    match, err := regexp.MatchString(urlRegeXp, email)
    if err != nil {
        panic(errMsgRegeXp)
    }
    return match
}


func IsURL(url string) bool {
    urlRegeXp := "^((http|https|ftp)\\://)?([a-zA-Z0-9\\.\\-]+(\\:[a-zA-Z0-9\\.&amp;\\$\\-]+)*@)?((25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9])|([a-zA-Z0-9\\-]+\\.)*[a-zA-Z0-9\\-]+\\.[a-zA-Z]{2,4})(\\:[0-9]+)?(/[^/][a-zA-Z0-9\\.\\,\\?\\'\\/\\+&amp;\\$#\\=~_\\-@]*)*$"
    match, err := regexp.MatchString(urlRegeXp, url)
    if err != nil {
        panic(err)
    }
    return match
}

func ResponseToString(resp *http.Response) (string, error) {
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        errMsg := "Can't Resolve Http Response. ioutil error: " + err.Error()
        fmt.Errorf(errMsg)
        return "", errors.New(errMsg)
    } else {
        return string(body), nil
    }
}

func HttpGet(url, username, password string) (*http.Response, error) {
    return Http(HTTP_METHOD_GET, url, username, password, nil)
}

func HttpPost(url, username, password string, data []byte) (*http.Response, error) {
    return Http(HTTP_METHOD_POST, url, username, password, data)
}

func HttpDelete(url, username, password string, data []byte) (*http.Response, error) {
    return Http(HTTP_METHOD_DELETE, url, username, password, data)
}

func HttpPut(url, username, password string, data []byte) (*http.Response, error) {
    return Http(HTTP_METHOD_PUT, url, username, password, data)
}

func IsValidHttpMethod(method string) bool {
    return strings.EqualFold(method, HTTP_METHOD_GET) || strings.EqualFold(method, HTTP_METHOD_POST) || strings.EqualFold(method, HTTP_METHOD_PUT) || strings.EqualFold(method, HTTP_METHOD_DELETE)
}

func Http(method, url, username, password string, data []byte) (*http.Response, error) {
    if !IsValidHttpMethod(method) {
        errMsg := "fail to send http request, parameter <method> is illegal."
        fmt.Errorf(errMsg)
        return nil, errors.New(errMsg)
    }

    if url == "" {
        errMsg := "fail to send http request, parameter <url> is empty."
        fmt.Errorf(errMsg)
        return nil, errors.New(errMsg)
    }

    var b io.Reader = nil

    if data != nil {
        b = bytes.NewBuffer(data)
    }

    client := &http.Client{}
    req, err := http.NewRequest(method, url, b)
    if err != nil {
        errMsg := "fail to send http request, " + err.Error()
        fmt.Errorf(errMsg)
        return nil, errors.New(errMsg)
    }
    req.Header.Add("Content-Type", "application/json")
    if username != "" && password != "" {
        req.SetBasicAuth(username, password)
    }
    resp, err := client.Do(req)
    if err != nil {
        errMsg := "fail to send http request, " + err.Error()
        fmt.Errorf(errMsg)
        return nil, errors.New(errMsg)
    }
    return resp, nil
}
