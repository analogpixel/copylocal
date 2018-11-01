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

  /*
  http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "Hello world!")
  })
  */
    http.HandleFunc("/upload/", ReceiveFile)

  //http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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
    
    // Copy the file data to my buffer
    io.Copy(&Buf, file)
    // do something with the contents...
    // I normally have a struct defined and unmarshal into a struct, but this will
    // work as an example
    contents := Buf.String()
    //fmt.Println(contents)
    ioutil.WriteFile(header.Filename, []byte(contents) , 0644)
    // I reset the buffer in case I want to use it again
    // reduces memory allocations in more intense projects
    Buf.Reset()
    // do something else
    // etc write header
    fmt.Fprintln(w, "Hello world!")
    return
}

func main() {
  var ip net.IP
  var dns_string string

  // Start out http service
  go startService()
  //log.SetFlags(0)
  //log.SetOutput(ioutil.Discard)
 
  ip = GetOutboundIP()
  dns_string = fmt.Sprintf("copy.local 60 IN A %s", ip)
  fmt.Println(dns_string)
  //mdns.Publish("copy.local 60 IN A 192.168.0.118")
  mdns.Publish(dns_string)

  // Sleep forever
  select{}
}

