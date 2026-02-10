package main

import (
	"fmt"
	"httpServer_go/internal/request"
	"log"
	"net"
)

// func getLinesChannel(f io.ReadCloser) <-chan string { //returns recieve only channel of strings
// 	out := make(chan string, 1)
// 	go func() {
// 		defer f.Close()  //close the function
// 		defer close(out) //close the channel
// 		str := ""
// 		for {
// 			data := make([]byte, 8)
// 			n, err := f.Read(data)
// 			if err != nil {
// 				break
// 			}
// 			i := bytes.IndexByte(data[:n], '\n')
// 			if i != -1 { //newline found, add all data till i to
// 				str += string(data[:i])
// 				data = data[i+1:]
// 				out <- str //instead of printing send out the str
// 				str = ""
// 			} else { //newline not found,
// 				str += string(data)
// 			}
// 			if len(str) != 0 {
// 				out <- str
// 			}
// 		}
// 	}()
// 	return out
// }

func main() {
	// file, err := os.Open("messages.txt")
	// if err != nil {
	// 	log.Fatal("error:", err)
	// }
	// defer file.Close()
	// lines:=getLinesChannel(file)

	listener, err := net.Listen("tcp", ":42069") //setting up a TCP listener on port 42069
	if err != nil {
		log.Fatal("error:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error:", err)
		}
		fmt.Printf("tcp connection accepted")
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error:", err)
		}
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method:%s\n", r.RequestLine.Method)
		fmt.Printf("- Target:%s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version:%s\n", r.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		r.Headers.ForEach(func(n, v string) {
			fmt.Printf("- %s: %s\n", n, v)
		})
		fmt.Printf("Body:\n")
		fmt.Printf("%s\n", r.Body)

	}
}
