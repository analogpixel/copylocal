package main

import (
  "fmt"
  "log"
  "net/http"
  //"os"
  "github.com/davecheney/mdns"
  //"io/ioutil"
  "bytes"
  "io"
  "strings"
)

// Our fake service.
// This could be a HTTP/TCP service or whatever you want.
func startService() {

  http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "Hello world!")
  })

    http.HandleFunc("/upload/", ReceiveFile)

  //http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  log.Println("starting http service...")
  if err := http.ListenAndServe(":80", nil); err != nil {
    log.Fatal(err)
  }
}

func ReceiveFile(w http.ResponseWriter, r *http.Request) {
    log.Print("Got a file now what")
    var Buf bytes.Buffer
    // in your case file would be fileupload
    file, header, err := r.FormFile("fileupload")
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
    fmt.Println(contents)
    // I reset the buffer in case I want to use it again
    // reduces memory allocations in more intense projects
    Buf.Reset()
    // do something else
    // etc write header
    fmt.Fprintln(w, "Hello world!")
    return
}

func main() {
  // Start out http service
  go startService()
  //log.SetFlags(0)
  //log.SetOutput(ioutil.Discard)
  // Setup our service export
  /*
  host, _ := os.Hostname()
  info := []string{"My awesome service"}
  service, _ := mdns.NewMDNSService(host, "_copyme._tcp", "", "", 80, nil, info)

  // Create the mDNS server, defer shutdown
  server, _ := mdns.NewServer(&mdns.Config{Zone: service})

  defer server.Shutdown()
  */
  mdns.Publish("copy.local 60 IN A 10.137.82.56")

  // Sleep forever
  select{}
}

