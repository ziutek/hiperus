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
		fmt.Fprintf(os.Stderr, "Usage: %s PWFILE\n", os.Args[0])
		os.Exit(1)
	}

	url := "https://backend.hiperus.pl:8080/hiperusapi.php"
	user, passwd, domain := readLoginInfo(os.Args[1])

	s, err := hiperus.NewSession(url, user, passwd, domain)
	checkErr(err)

	/*c := &hiperus.Customer{
		Name:            "test4",
		Email:           "test4@Lnet.pl",
		Address:         "Służbowa",
		StreetNumber:    "135",
		ApartmentNumber: "162",
		PostCode:        "92-305",
		City:            "Łódź",
		Country:         "Polska",

		BiName:            "test4",
		BiAddress:         "Służbowa",
		BiStreetNumber:    "135",
		BiApartmentNumber: "162",
		BiPostCode:        "92-305",
		BiCity:            "Łódź",
		BiCountry:         "Polska",
		BiNIP:             "777-222-44-11",
		BiRegon:           "123456",

		ExtBillingId: 9,

		PaymentType: "postpaid",
		IsWLR:       true,

		DefaultPriceListId:    1843,
		ConsentDataProcessing: true,
	}

	id, err := s.CreateCustomer(c)
	checkErr(err)
	fmt.Println("Utworzono uzytkownika o id:", id)*/

	var customer hiperus.Customer

	id := uint32(12982)
	fmt.Println("Klient o id:", id)
	checkErr(s.GetCustomerData(&customer, id))
	fmt.Printf("%+v\n", customer)

	id = 2
	fmt.Println("Klient o extId:", id)
	checkErr(s.GetCustomerDataExtId(&customer, id))
	fmt.Printf("%+v\n", customer)

	fmt.Println("Lista klientów:")
	cl, err := s.GetCustomerList(0, 0, "")
	checkErr(err)
	for cl.Next() {
		checkErr(cl.Scan(&customer))
		fmt.Printf("%+v\n", customer)
	}

	// Biling
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
	fmt.Println("\n\nBiling za okres od", start, "do", stop, "\n")
	var call hiperus.Call
	for b.Next() {
		checkErr(b.Scan(&call))
		fmt.Printf("%+v\n", call)
	}

	// Numery
	num, cc, err := s.GetFirstFreePlatformNumber("42")
	checkErr(err)
	fmt.Println("\n\nPierwszy wolny numer:", cc, num)

	// Cenniki
	pl, err := s.GetCustomerPricelistList()
	checkErr(err)
	fmt.Println("\n\nCenniki:")
	var p hiperus.CustomerPricelist
	for pl.Next() {
		checkErr(pl.Scan(&p))
		fmt.Printf("%+v\n", p)
	}

	id = 1843
	fmt.Println("\n\nCennik o id:", id)
	p, err = s.GetCustomerPricelist(id, "")
	checkErr(err)
	fmt.Printf("%+v\n", p)

	name := "cn1"
	fmt.Println("\n\nCennik o nazwie:", name)
	p, err = s.GetCustomerPricelist(0, name)
	checkErr(err)
	fmt.Printf("%+v\n", p)

	// Terminale
	tl, err := s.GetTerminalList(12372, 0, 0)
	checkErr(err)
	var term hiperus.Terminal
	for tl.Next() {
		checkErr(tl.Scan(&term))
		fmt.Printf("%+v\n", term)
	}
}
