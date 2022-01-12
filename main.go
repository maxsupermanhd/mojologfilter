package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

func main() {
	re := regexp.MustCompile(`(?m)Ignoring chunk since it's not in the view range: (-?[0-9]+), (-?[0-9]+)`)
	re2 := regexp.MustCompile(`(?m)\[[0-9]+:[0-9]+:[0-9]+] \[Render thread/INFO\]: \[CHAT\] [~a-zA-Z0-9_]{2,16} joined\.`)
	re3 := regexp.MustCompile(`(?m)\[[0-9]+:[0-9]+:[0-9]+] \[Render thread/INFO\]: \[CHAT\] [~a-zA-Z0-9_]{2,16} left\.`)
	entrycount := 0
	out, err := os.Create(os.Args[1])
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer out.Close()
	for _, fn := range os.Args[2:] {
		if fn == os.Args[1] {
			continue
		}
		f, err := os.Open(fn)
		if err != nil {
			log.Print(err.Error())
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			e := re.FindStringSubmatch(scanner.Text())
			if e != nil {
				entrycount++
				out.Write([]byte(e[0] + "\n"))
				continue
			}
			e2 := re2.FindStringSubmatch(scanner.Text())
			if e2 != nil {
				out.Write([]byte(e2[0] + "\n"))
				continue
			}
			e3 := re3.FindStringSubmatch(scanner.Text())
			if e3 != nil {
				out.Write([]byte(e3[0] + "\n"))
				continue
			}
		}
	}
}
