package hiperus

import (
	"errors"
	"github.com/ziutek/soap"
	"time"
)

type Customer struct {
	Id uint32 `soap:"id,in"` // identyfikator klienta

	Name            string `soap:"name"`             // nazwa klienta
	Email           string `soap:"email"`            // email klienta
	Address         string `soap:"address"`          // ulica/miejscowość
	StreetNumber    string `soap:"street_number"`    // numer budynku
	ApartmentNumber string `soap:"apartment_number"` // numer mieszkania
	PostCode        string `soap:"postcode"`         // kod pocztowy
	City            string `soap:"city"`             // miasto
	Country         string `soap:"country"`          // kraj

	// Dane do celow bilingowych (na fakture)
	BiName            string `soap:"b_name"`
	BiAddress         string `soap:"b_address"`
	BiStreetNumber    string `soap:"b_street_number"`
	BiApartmentNumber string `soap:"b_apartment_number"`
	BiPostCode        string `soap:"b_postcode"`
	BiCity            string `soap:"b_city"`
	BiCountry         string `soap:"b_country"`
	BiNIP             string `soap:"b_nip"`
	BiRegon           string `soap:"b_regon"`

	ExtBillingId uint32 `soap:"ext_billing_id"` // id z syst. zewnętrznego

	IssueInvoice bool   `soap:"issue_invoice"` // wystawiac faktury
	PaymentType  string `soap:"payment_type"`  // typ platn. {prepaid, postpaid}
	IsWLR        bool   `soap:"is_wlr"`        // czy klient WLR

	DefaultPriceListId    uint32 `soap:"id_default_pricelist"`
	ConsentDataProcessing bool   `soap:"consent_data_processing"`

	ResellerId              uint32 `soap:"id_reseller,in"`
	Active                  bool   `soap:"active,in"`
	OpenRegistration        bool   `soap:"open_registration,in"`
	IsRemoved               bool   `soap:"is_removed,in"`
	DefaultBalanceId        uint32 `soap:"id_default_balance,in"`
	CustomerPostpaidLimitId uint32 `soap:"id_customer_postpaid_limit,in"`

	CreateDate           time.Time `soap:"create_date,in"`
	PlatformUserAddStamp string    `soap:"platform_user_add_stamp,in"`
}

func firstRow(rs *soap.Element) (*soap.Element, error) {
	if len(rs.Children) == 0 {
		return nil, errors.New("hiperus: empty result set")
	}
	return rs.Children[0], nil
}

func (s *Session) CreateCustomer(c *Customer) (uint32, error) {
	rs, err := s.cmd("AddCustomer", c)
	if err != nil {
		return 0, err
	}
	e, err := firstRow(rs)
	if err != nil {
		return 0, err
	}
	if e, err = e.Get("id"); err != nil {
		return 0, err
	}
	return e.AsUint32()
}

type getCustomerData struct {
	CustomerId uint32 `soap:"id_customer"`
}

func (s *Session) GetCustomerData(c *Customer, id uint32) error {
	rs, err := s.cmd("GetCustomerData", getCustomerData{id})
	if err != nil {
		return err
	}
	e, err := firstRow(rs)
	if err != nil {
		return err
	}
	return e.LoadStruct(c, false)
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

func (b *CustomerList) Scan(c *Customer) (err error) {
	row := b.rs.Children[b.n-1]
	return row.LoadStruct(c, false)
}

type getCustomerList struct {
	Offset int     `soap:"offset"`
	Limit  *int    `soap:"limit"`
	Query  *string `soap:"query"`
}

// GetCustomerList Zwraca listę klientów utworzonych na platformie Hiperus C5
//	offset – od którego rekordu pobrać dane,
//	limit  - jak dużo danych pobrać (0 oznacza brak limitu),
//	query  - pobrać tylko klientów których nazwa rozpoczyna sie od query,
//	         (wielkość znaków nie ma znaczenia, "" pasuje do wszystkich)
func (s *Session) GetCustomerList(offset, limit int, query string) (
	*CustomerList, error) {

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
	return &CustomerList{rs, 0}, nil
}
