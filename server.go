package main

import (
  "io"
  "net/http"
  "context"
  "time"
)

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
  w.Write(body)
}
