package main

import (
	"bufio"
	"fmt"
	"github.com/ziutek/hiperus"
	"os"
	"time"
)

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func readLoginInfo(fname string) (user, passwd, domain string) {
	pwf, err := os.Open(fname)
	checkErr(err)
	defer pwf.Close()

	s := bufio.NewScanner(pwf)

	if s.Scan() {
		user = s.Text()
	}
	if s.Scan() {
		passwd = s.Text()
	}
	if s.Scan() {
		domain = s.Text()
	}
	checkErr(s.Err())
	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s PWDFILE\n", os.Args[0])
		os.Exit(1)
	}

	url := "https://backend.hiperus.pl:8080/hiperusapi.php"
	user, passwd, domain := readLoginInfo(os.Args[1])

	s, err := hiperus.NewSession(url, user, passwd, domain)
	checkErr(err)

	// Billing
	start := time.Date(
		2013, 4, 16,
		23, 59, 59, 0,
		time.Local,
	)
	stop := start.Add(30 * 24 * time.Hour)
	b, err := s.GetBilling(
		start, stop,
		0, 0,
		true, 0, "outgoing",
	)
	checkErr(err)
	fmt.Println("Biling za okres od", start, "do", stop)
	var c hiperus.Call
	for b.Next() {
		checkErr(b.Scan(&c))
		fmt.Printf("%+v\n", c)
	}

}
