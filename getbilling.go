package hiperus

import (
	"github.com/ziutek/soap"
	"time"
)

type Call struct {
	Id uint64

	Start    time.Time     // moment nawiazania polaczenia (time.Local)
	Duration time.Duration // czas trwania polczaenia
	RelCause ReleaseCode   // powód zakończenia (ISDN ISUP release code)

	CustomerName string // nazwa klienta
	TerminalName string // nazwa terminala
	ExtBillingId uint32 // id klienta w systemie lokalnym
	Caller       string // numer, z którego dzwoniono, format krajowy
	BillCPB      string // numer, na który dzwoniono format międzynarodowy bez +
	CallType     string // {incoming, outgoing, disa, forwarded, internal, vpbx}

	Country     string // kierunek - kraj
	Description string // opis kierunku
	Operator    string // operator kierunku
	Type        string // {mobile, geographic, premium, aus, NGN}

	Price      float64 // cena połączenia dla klienta
	Cost       float64 // koszt połączenia dla klienta
	InitCharge float64 // opłata za inicjacje połączenia dla klienta

	RePrice          float64 // cena połączenia dla resellera
	ReCost           float64 // koszt połączenia dla resellera
	ReInitCharge     float64 // opłata za inicjację połączenia dla resellera
	Margin           float64 // marża resellera
	SubscriptionUsed bool    // czy został użyty abonament (np. darmowe minuty)

	PlatformType string // typ usługi platformy Hiperus C5 {VOIP, WLR}
}

type Billing struct {
	rs *soap.Element
	n  int
}

func (b *Billing) Next() bool {
	if b.n == len(b.rs.Children) {
		return false
	}
	b.n++
	return true
}

func (b *Billing) Scan(c *Call) (err error) {
	row := b.rs.Children[b.n-1]
	var (
		e   *soap.Element
		u8  byte
		u32 uint32
	)

	if e, err = row.Get("id"); err != nil {
		return
	}
	if c.Id, err = e.AsUint64(); err != nil {
		return
	}

	if e, err = row.Get("start_time"); err != nil {
		return
	}
	if c.Start, err = e.AsTime(); err != nil {
		return
	}

	if e, err = row.Get("duration"); err != nil {
		return
	}
	if u32, err = e.AsUint32(); err != nil {
		return
	}
	c.Duration = time.Duration(u32) * time.Second

	if e, err = row.Get("rel_cause"); err != nil {
		return
	}
	if u8, err = e.AsUint8(); err != nil {
		return
	}
	c.RelCause = ReleaseCode(u8)

	if e, err = row.Get("customer_name"); err != nil {
		return
	}
	if c.CustomerName, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("terminal_name"); err != nil {
		return
	}
	if c.TerminalName, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("ext_billing_id"); err != nil {
		return
	}
	if e.Nil {
		c.ExtBillingId = 0
	} else if c.ExtBillingId, err = e.AsUint32(); err != nil {
		return
	}

	if e, err = row.Get("caller"); err != nil {
		return
	}
	if c.Caller, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("bill_cpb"); err != nil {
		return
	}
	if c.BillCPB, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("calltype"); err != nil {
		return
	}
	if c.CallType, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("country"); err != nil {
		return
	}
	if e.Nil {
		c.Country = ""
	} else if c.Country, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("description"); err != nil {
		return
	}
	if e.Nil {
		c.Description = ""
	} else if c.Description, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("operator"); err != nil {
		return
	}
	if e.Nil {
		c.Operator = ""
	} else if c.Operator, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("type"); err != nil {
		return
	}
	if e.Nil {
		c.Type = ""
	} else if c.Type, err = e.Str(); err != nil {
		return
	}

	if e, err = row.Get("price"); err != nil {
		return
	}
	if c.Price, err = e.AsFloat64(); err != nil {
		return
	}

	if e, err = row.Get("cost"); err != nil {
		return
	}
	if c.Cost, err = e.AsFloat64(); err != nil {
		return
	}

	if e, err = row.Get("init_charge"); err != nil {
		return
	}
	if c.InitCharge, err = e.AsFloat64(); err != nil {
		return
	}

	if e, err = row.Get("reseller_price"); err != nil {
		return
	}
	if c.RePrice, err = e.AsFloat64(); err != nil {
		return
	}

	if e, err = row.Get("reseller_cost"); err != nil {
		return
	}
	if c.ReCost, err = e.AsFloat64(); err != nil {
		return
	}

	if e, err = row.Get("reseller_init_charge"); err != nil {
		return
	}
	if c.ReInitCharge, err = e.AsFloat64(); err != nil {
		return
	}

	if e, err = row.Get("margin"); err != nil {
		return
	}
	if c.Margin, err = e.AsFloat64(); err != nil {
		return
	}

	if e, err = row.Get("subscription_used"); err != nil {
		return
	}
	c.SubscriptionUsed = asBool(e)

	if e, err = row.Get("platform_type"); err != nil {
		return
	}
	if c.PlatformType, err = e.Str(); err != nil {
		return
	}

	return nil
}

type getBilling struct {
	From         time.Time `xml:"from"`
	To           time.Time `xml:"to"`
	Offset       int       `xml:"offset"`
	Limit        *int      `xml:"limit"`
	Compress     bool      `xml:"compress"`
	SuccessCalls bool      `xml:"success_calls"`
	CallType     *string   `xml:"calltype"`
	CustomerId   *int      `xml:"id_customer"`
}

// GetBilling pobiera dane billingowe z platformy Hiperus C5
func (s *Session) GetBilling(from, to time.Time, offset, limit int,
	successCalls bool, customerId int, callType string) (*Billing, error) {

	arg := &getBilling{
		From:         from,
		To:           to,
		Offset:       offset,
		SuccessCalls: successCalls,
	}
	if limit != 0 {
		arg.Limit = &limit
	}
	if callType != "" {
		arg.CallType = &callType
	}
	if customerId != 0 {
		arg.CustomerId = &customerId
	}
	rs, err := s.cmd("GetBilling", arg)
	if err != nil {
		return nil, err
	}
	return &Billing{rs, 0}, nil
}
