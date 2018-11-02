package main

import (
  "fmt"
  "log"
  "net/http"
  //"os"
  "github.com/davecheney/mdns"
  "io/ioutil"
  "bytes"
  "io"
  "strings"
  "net"
  "github.com/gobuffalo/packr"
)

// Our fake service.
// This could be a HTTP/TCP service or whatever you want.
func startService() {

  box := packr.NewBox("./templates")

  http.Handle("/", http.FileServer(box))
  http.HandleFunc("/upload/", ReceiveFile)

  log.Println("starting http service...")
  if err := http.ListenAndServe(":80", nil); err != nil {
    log.Fatal(err)
  }
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}

func ReceiveFile(w http.ResponseWriter, r *http.Request) {
    log.Print("Got a file now what")
    var Buf bytes.Buffer
    file, header, err := r.FormFile("file")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    name := strings.Split(header.Filename, ".")
    fmt.Printf("File name %s\n", name[0])
    fmt.Printf("Full name %s\n", header.Filename)
    io.Copy(&Buf, file)
    contents := Buf.String()
    ioutil.WriteFile(header.Filename, []byte(contents) , 0644)
    Buf.Reset()
    return
}

func main() {
  var ip net.IP
  var dns_string string

  // Start out http service
  go startService() 
  ip = GetOutboundIP()
  dns_string = fmt.Sprintf("copy.local 60 IN A %s", ip)
  fmt.Println(dns_string)
  mdns.Publish(dns_string)

  // Sleep forever
  select{}
}

