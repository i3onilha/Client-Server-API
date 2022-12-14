package main

import (
  "io"
  "net/http"
  "database/sql"
  "encoding/json"
  "context"
  "time"
  "fmt"
  _ "github.com/mattn/go-sqlite3"
)

type DollarExchange struct {
  Usdbrl `json:"USDBRL"`
}

type Usdbrl struct {
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
}

func init() {
  fmt.Println(createDollarPriceTable())
}

func main() {
  http.HandleFunc("/cotacao", dollarExchange)
  http.ListenAndServe(":8090", nil)
}

func dollarExchange(w http.ResponseWriter, r *http.Request) {
  ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
  defer cancel()
  urlApi := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
  req, err := http.NewRequestWithContext(ctx, "GET", urlApi, nil)
  if err != nil {
    w.Write([]byte(err.Error()))
    return
  }
  res, err := http.DefaultClient.Do(req)
  if err != nil {
    w.Write([]byte(err.Error()))
    return
  }
  defer res.Body.Close()
  body, err := io.ReadAll(res.Body)
  if err != nil {
    w.Write([]byte(err.Error()))
    return
  }
  var ag DollarExchange
  json.Unmarshal(body, &ag)
  err = storeDollarPrice(ag.Usdbrl)
  if err != nil {
    w.Write([]byte(err.Error()))
    return
  }
  w.Write(body)
}

func storeDollarPrice(p Usdbrl) error {
  db, err := sql.Open("sqlite3", "./dollar_price.db")
  if err != nil {
    return err
  }
  rawSQl := "insert into dollar_price (code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date)"
  rawSQl += " values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
  stmt, err := db.Prepare(rawSQl)
  if err != nil {
    return err
  }
  stmt.Exec(p.Code, p.Codein, p.Name, p.High, p.Low, p.VarBid, p.PctChange, p.Bid, p.Ask, p.Timestamp, p.CreateDate)
  if err != nil {
    return err
  }
  return nil
}

func createDollarPriceTable() error {
  db, err := sql.Open("sqlite3", "./dollar_price.db")
  if err != nil {
    return err
  }
  rawSQl := `create table
    dollar_price (
      code text,
      codein text,
      name text,
      high text,
      low text,
      var_bid text,
      pct_change text,
      bid text,
      ask text,
      timestamp text,
      create_date text
    );
  `
  stmt, err := db.Prepare(rawSQl)
  if err != nil {
    return err
  }
  stmt.Exec()
  if err != nil {
    return err
  }
  return nil
}
