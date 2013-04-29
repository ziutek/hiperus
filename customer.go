package hiperus

import (
	"github.com/ziutek/soap"
	"time"
)

type Customer struct {
	Id uint32 `xml:"-"` // identyfikator klienta

	Name         string `xml:"name"`          // nazwa klienta
	Email        string `xml:"email"`         // email klienta
	Address      string `xml:"address"`       // adres linia1: ulica miejscowość
	StreetNumber string `xml:"street_number"` // numer budynku
	FlatNumber   string `xml:"flat_number"`   // numer mieszkania
	PostCode     string `xml:"postcode"`      // kod pocztowy
	City         string `xml:"city"`          // miasto
	Country      string `xml:"country"`       // kraj

	// Dane do celow bilingowych (na fakture)
	BiName         string `xml:"b_name"`
	BiAddress      string `xml:"b_address"`
	BiStreetNumber string `xml:"b_street_number"`
	BiFlatNumber   string `xml:"b_flat_number"`
	BiPostcode     string `xml:"b_postcode"`
	BiCity         string `xml:"b_city"`
	BiCountry      string `xml:"b_country"`
	BiNIP          string `xml:"b_nip"`
	BiRegon        string `xml:"b_regon"`

	ExtBillingId uint32 `xml:"ext_billing_id"` // id z syst. zewnętrznego

	IssueInvoice bool   `xml:"issue_invoice"` // wystawiac faktury
	PaymentType  string `xml:"payment_type"`  // typ platn. {prepaid, postpaid}
	IsWLR        bool   `xml:"is_wlr"`        // czy klient WLR

	ConsentDataProcessing bool `xml:"consent_data_processing"`

	DefaultPriceListId      uint32    `xml:"-"` // id domyślnego cennika
	ResellerId              uint32    `xml:"-"` // id resellera
	Active                  bool      `xml:"-"` // czy klient aktywny
	CreateDate              time.Time `xml:"-"` // data utworzenia klienta
	PlatformUserAddStamp    string    `xml:"-"` // kto utworzył klienta
	OpenRegistration        bool      `xml:"-"` // utw. przez otwartą rejestr.
	IsRemoved               bool      `xml:"-"` // czy usunięty
	DefaultBalanceId        uint32    `xml:"-"` // id. domyślnego limitu prepaid
	CustomerPostpaidLimitId uint32    `xml:"-"`
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
	var e *soap.Element

	if e, err = row.Get("id"); err != nil {
		return
	}
	if c.Id, err = e.AsUint32(); err != nil {
		return
	}

	if e, err = row.Get("name"); err != nil {
		return
	}
	if c.Name, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("email"); err != nil {
		return
	}
	if c.Email, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("address"); err != nil {
		return
	}
	if c.Address, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("street_number"); err != nil {
		return
	}
	if c.StreetNumber, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("apartment_number"); err != nil {
		return
	}
	if c.FlatNumber, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("postcode"); err != nil {
		return
	}
	if c.PostCode, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("city"); err != nil {
		return
	}
	if c.City, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("country"); err != nil {
		return
	}
	if c.Country, err = e.Str(); err != nil {
		return
	}
	if e, err = row.Get("b_name"); err != nil {
		return
	}
	if c.Name, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("b_address"); err != nil {
		return
	}
	if c.Address, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("b_street_number"); err != nil {
		return
	}
	if c.StreetNumber, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("b_apartment_number"); err != nil {
		return
	}
	if c.FlatNumber, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("b_postcode"); err != nil {
		return
	}
	if c.PostCode, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("b_city"); err != nil {
		return
	}
	if c.City, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("b_country"); err != nil {
		return
	}
	if c.Country, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("b_nip"); err != nil {
		return
	}
	if c.BiNIP, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("b_regon"); err != nil {
		return
	}
	if c.BiRegon, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("ext_billing_id"); err != nil {
		return
	}
	if c.ExtBillingId, err = e.AsUint32(); err != nil {
		return
	}

	if e, err = row.Get("issue_invoice"); err != nil {
		return
	}
	c.IssueInvoice = asBool(e)

	if e, err = row.Get("payment_type"); err != nil {
		return
	}
	if c.PaymentType, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("is_wlr"); err != nil {
		return
	}
	c.IsWLR = asBool(e)

	if e, err = row.Get("consent_data_processing"); err != nil {
		return
	}
	c.ConsentDataProcessing = asBool(e)

	if e, err = row.Get("id_default_pricelist"); err != nil {
		return
	}
	if c.DefaultPriceListId, err = e.AsUint32(); err != nil {
		return
	}

	if e, err = row.Get("id_reseller"); err != nil {
		return
	}
	if c.ResellerId, err = e.AsUint32(); err != nil {
		return
	}

	if e, err = row.Get("active"); err != nil {
		return
	}
	c.Active = asBool(e)

	if e, err = row.Get("create_date"); err != nil {
		return
	}
	if c.CreateDate, err = e.AsTime(); err != nil {
		return
	}

	if e, err = row.Get("platform_user_add_stamp"); err != nil {
		return
	}
	if c.PlatformUserAddStamp, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("open_registration"); err != nil {
		return
	}
	c.OpenRegistration = asBool(e)

	if e, err = row.Get("is_removed"); err != nil {
		return
	}
	c.IsRemoved = asBool(e)

	if e, err = row.Get("id_default_balance"); err != nil {
		return
	}
	if e.Nil {
		c.DefaultBalanceId = 0
	} else if c.DefaultBalanceId, err = e.AsUint32(); err != nil {
		return
	}

	if e, err = row.Get("id_customer_postpaid_limit"); err != nil {
		return
	}
	if e.Nil {
		c.CustomerPostpaidLimitId = 0
	} else if c.CustomerPostpaidLimitId, err = e.AsUint32(); err != nil {
		return
	}

	return
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
