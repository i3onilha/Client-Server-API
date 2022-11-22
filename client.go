package main


import (
  "io"
  "os"
  "net/http"
  "encoding/json"
  "time"
  "fmt"
  "errors"
)

type DollarExchange struct {
  Usdbrl struct {
    Code       string `json:"code"`
    Codein     string `json:"codein"`
    Name       string `json:"name"`
    High       string `json:"high"`
    Low        string `json:"low"`
    VarBid     string `json:"varBid"`
    PctChange  string `json:"pctChange"`
    Bid        string `json:"bid"`
    Ask        string `json:"ask"`
    Timestamp  string `json:"timestamp"`
    CreateDate string `json:"create_date"`
  } `json:"USDBRL"`
}

const urlApi string = "http://10.58.64.197:8090/cotacao"

func main() {
  c := http.Client{Timeout: time.Millisecond*300}
  res, err := c.Get(urlApi)
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  defer res.Body.Close()
  body, err := io.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  var ag DollarExchange
  json.Unmarshal(body, &ag)
  if ag.Usdbrl.Bid == "" {
    fmt.Println(errors.New("Bid empty"))
    return
  }
  f, err := os.OpenFile("cotacao.txt", os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  defer f.Close()
  price := ag.Usdbrl.Bid + "\n"
  if _, err := f.Write([]byte(price)); err != nil {
    fmt.Println(err.Error())
  }
}
