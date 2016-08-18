package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

func read(s *serial.Port, format string) {
	buf := make([]byte, 128)
	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		switch format {
		case "h":
			for _, v := range buf[:n] {
				fmt.Printf("%0#2x ", v)
			}
		case "d":
			for _, v := range buf[:n] {
				fmt.Print(v, " ")
			}
		case "a":
			fmt.Print(string(buf[:n]))
		}

	}

}
func write(s *serial.Port) {
	reader := bufio.NewReader(os.Stdin)
	for {
		bytes, err := reader.ReadBytes(0)
		if err != nil {
			log.Fatal(err)
		}
		s.Write(bytes)
	}
}
func main() {
	port := flag.String("port", "COM4", "")
	flag.StringVar(port, "p", "COM4", "port to connect to i.e. COM4")

	speed := flag.Uint("baud", 57600, "")
	flag.UintVar(speed, "b", 57600, "baudrate")
	format := flag.String("format", "a", "")
	flag.StringVar(format, "f", "h", "set format h: hex, d decimal, a: ascii")

	flag.Parse()
	trimmed := strings.TrimPrefix(*port, "COM")
	_, err := strconv.ParseInt(trimmed, 10, 32)
	if !strings.HasPrefix(*port, "COM") || err != nil {
		fmt.Println("error in port value provided")
		flag.PrintDefaults()
		os.Exit(1)
	}
	c := &serial.Config{Name: *port, Baud: int(*speed)}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	go read(s, *format)
	write(s)

}
