package main

import (
	"io"
	"log"
	"net"
	"sync"
	"github.com/Trym123/is105sem03/mycrypt"
	"github.com/Trym123/minyr/yr"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.2:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					switch msg := string(dekryptertMelding); msg {
  				        case "ping":
						//_, err = c.Write([]byte("pong"))
						krypterMelding := mycrypt.Krypter([]rune(string("pong")), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
						log.Println("krypter melding: ", string(krypterMelding))
						_, err = conn.Write([]byte(string(krypterMelding)))
					case "Kjevik;SN39040;18.03.2022 01:50;6":
						newString, err := yr.CelsiusToFahrenheitLine("Kjevik;SN39040;18.03.2022 01:50;6")
						if err != nil {
							log.Fatal(err)
						}
							//dividedString := strings.Split("Kjevik;SN39040;18.03.2022 01:50;6", ";")
							
							//if fahr, err := strconv.ParseFloat(dividedString[3], 64); err == nil {
							//log.Println(conv.CelsiusToFarhenheit(fahr)) }
							//joinedString := strings.Join(dividedString, ";")
						_, err = conn.Write([]byte(string(newString)))
					default:
						_, err = c.Write(buf[:n])
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				
				}
			}(conn)
		}
	}()
	wg.Wait()
}
