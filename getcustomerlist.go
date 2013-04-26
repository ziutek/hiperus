package hiperus

import (
	"github.com/ziutek/soap"
	"time"
)

type Customer struct {
	Id uint64 // identyfikator klienta

	Name            string // nazwa klienta
	Email           string // email klienta
	Address         string // adres linia1 – nazwa ulicy / miejscowość
	StreetNumber    string // numer ulicy / miejscowości
	ApartmentNumber string // numer mieszkania
	PostCode        string // kod pocztowy
	City            string // miasto
	Country         string // kraj

	// Dane na fakture
	BiName            string
	BiAddress         string
	BiStreetNumber    string
	BiApartmentNumber string
	BiPostcode        string
	BiCity            string
	BiCountry         string
	BiNIP             string
	BiRegon           string

	ResellerId   uint32 // identyfikator resellera
	ExtBillingId uint32 // identyfkator klienta z systemu zewnętrznego

	IssueInvoice            bool      // czy wystawiane faktury
	DefaultPricelistId      uint32    // identyfikator domyślnego cennika
	PaymentType             string    // typ płatności {prepaid, postpaid}
	DefaultBalanceId        string    // identyfikator domyślnego limitu prepaid
	Active                  bool      // czy klient aktywny
	IsWLR                   bool      // czy klient WLR
	ConsentDataProcessing   bool      // czy zgoda na przetwarzanie danych
	CreateDate              time.Time // data i czas utworzenia rekordu klienta
	PlatformUserAddStamp    string    // użytkownik który utworzył klienta
	OpenRegistration        bool      // czy utworzony poprzez otwartą rejestracje
	IsRemoved               bool      //  czy usunięty (dla Get... zawsze false)
	CustomerPostpaidLimitId uint32
}

type CustomerList struct {
	rs *soap.Element
	n  int
}

func (cl *CustomerList) Next() bool {
	if cl.n == len(cl.rs.Children) {
		return false
	}
	cl.n++
	return true
}

type getCustomerList struct {
	Offset int     `xml:"offset"`
	Limit  *int    `xml:"limit"`
	Query  *string `xml:"query"`
}

// GetCustomerList Zwraca listę klientów utworzonych na platformie Hiperus C5
//	offset – od którego rekordu pobrać dane,
//	limit  - jak dużo danych pobrać (0 oznacza brak limitu),
//	query  - pobrać tylko klientów których nazwa rozpoczyna sie od query,
//	         (wielkość znaków nie ma znaczenia, "" pasuje do wszystkich)
func (s *Session) GetCustomerList(offset, limit int, query string) (
	*soap.Element, error) {

	arg := getCustomerList{
		Offset: offset,
	}
	if limit != 0 {
		arg.Limit = &limit
	}
	if query != "" {
		arg.Query = &query
	}
	rs, err := s.cmd("GetCustomerList", arg)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
