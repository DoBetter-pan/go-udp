/** 
* @file udpclient.go
* @brief udp client
* @author yingx
* @date 2016-01-22
*/
package main

import (
    "fmt"
    "flag"
    "net"
    "bytes"
    "encoding/gob"
)

type params struct {
    host string
    port int
    enc string
    maintype string
    subtype  string
    message string
}

type event struct {
    Maintype string
    Subtype  string
    Message  string
}

func (data *event) encode(enc string) (string, error) {
    if enc == "gob" {
        var buffer bytes.Buffer
        enc := gob.NewEncoder(&buffer)
        err := enc.Encode(data)
        return buffer.String(), err
    } else if enc == "plain" {
        str := fmt.Sprintf("%s:%s:%s", data.Maintype, data.Subtype, data.Message)
        return str, nil
    } else {
        return "", fmt.Errorf("%s", "Not supported encoding way.") 
    }
}

func handleCommandLine() (*params, []string ){
    p := params{}

    flag.StringVar(&p.host, "host", "127.0.0.1", "host ot send udp to")
    flag.IntVar(&p.port, "port", 9898, "port ot send udp to")
    flag.StringVar(&p.enc, "enc", "plain", "the way of encoding: gob, plain")
    flag.StringVar(&p.maintype, "maintype", "log", "main type")
    flag.StringVar(&p.subtype, "subtype", "error", "sub type")
    flag.StringVar(&p.message, "message", "hello world", "message")
    flag.Parse()

    messages := flag.Args()

    return &p, messages
}

func main() {
    p, messages := handleCommandLine()
    ip := net.ParseIP(p.host)
    srcAddr := &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 0}
    dstAddr := &net.UDPAddr{IP: ip, Port: p.port}

    conn, err := net.ListenUDP("udp", srcAddr)
    //conn, err := net.DialUDP("udp", srcAddr, dstAddr)
    if err != nil {
        fmt.Println("failed to dial udp server:", err)
        return
    }
    defer conn.Close()

    for _, msg := range messages {
        encoder :=  event{p.maintype, p.subtype, msg}
        message, err := encoder.encode(p.enc)
        if err == nil {
            n, err := conn.WriteToUDP([]byte(message), dstAddr)
            if err != nil {
                fmt.Println("failed to send udp server:", err)
            } else {
                fmt.Printf("send %d data to udp server:%s\n", n, message)
            }
        } else {
            fmt.Println("failed to encode the data:", err)
        }
    }
}

