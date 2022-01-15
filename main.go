package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	re0 = regexp.MustCompile(`(?m)\[([0-9])+:([0-9]+):([0-9]+)\] \[Render thread/WARN\]: Ignoring chunk since it's not in the view range: (-?[0-9]+), (-?[0-9]+)`)
	re1 = regexp.MustCompile(`(?m)\[([0-9])+:([0-9]+):([0-9]+)\] \[Render thread/INFO\]: \[CHAT\] ([~a-zA-Z0-9_]{2,16}) joined\.`)
	re2 = regexp.MustCompile(`(?m)\[([0-9])+:([0-9]+):([0-9]+)\] \[Render thread/INFO\]: \[CHAT\] ([~a-zA-Z0-9_]{2,16}) left\.`)
	re3 = regexp.MustCompile(`(?m)([0-9]+)-([0-9]+)-([0-9]+)-([0-9]+)\.log`)
	re4 = regexp.MustCompile(`(?m)\[([0-9])+:([0-9]+):([0-9]+)\] \[Render thread/INFO\]: \[CHAT\] Setting user: ([~a-zA-Z0-9_]{2,16})`)
)

func main() {
	entrycount := 0
	logcount := 0
	out, err := os.Create(os.Args[1] + ".mlcd")
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer out.Close()
	out.Write([]byte("MLCD1\n"))
	for _, fn := range os.Args[2:] {
		if fn == os.Args[1] {
			continue
		}
		if !strings.HasSuffix(fn, ".log") {
			log.Printf("Filename [%s] does not end with .log, skipping.", fn)
			continue
		}
		filedate := "00000000"
		fnre := re3.FindStringSubmatch(fn)
		if fnre == nil {
			log.Printf("Filename [%s] does not match re3! Contents will not have timings information!", fn)
		} else {
			filedate = fnre[1]
			filedate += fnre[2]
			filedate += fnre[3]
		}
		f, err := os.Open(fn)
		if err != nil {
			log.Printf("%s, skipping.", err.Error())
			continue
		}
		defer f.Close()
		logcount++
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			e := re0.FindStringSubmatch(scanner.Text())
			if e != nil {
				entrycount++
				towrite := filedate + e[1] + e[2] + e[3] + " c " + e[4] + " " + e[5] + "\n"
				out.Write([]byte(towrite))
				continue
			}
			e = re1.FindStringSubmatch(scanner.Text())
			if e != nil {
				towrite := filedate + e[1] + e[2] + e[3] + " j " + e[4] + "\n"
				out.Write([]byte(towrite))
				continue
			}
			e = re2.FindStringSubmatch(scanner.Text())
			if e != nil {
				towrite := filedate + e[1] + e[2] + e[3] + " l " + e[4] + "\n"
				out.Write([]byte(towrite))
				continue
			}
			e = re4.FindStringSubmatch(scanner.Text())
			if e != nil {
				towrite := filedate + e[1] + e[2] + e[3] + " p " + e[4] + "\n"
				out.Write([]byte(towrite))
				continue
			}
		}
	}
	log.Printf("Collected %d coordinates from %d log file(s).", entrycount, logcount)
}
