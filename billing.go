package hiperus

import (
	"github.com/ziutek/soap"
	"time"
)

type Call struct {
	Id uint64 `soap:"id"`

	Start    time.Time     `soap:"start_time"` // moment nawiązania połączenia (localtime)
	Duration time.Duration `soap:"duration"`   // czas trwania połączenia
	RelCause ReleaseCode   `soap:"rel_cause"`  // powód zakończenia połączenia

	CustomerName string `soap:"customer_name"`  // nazwa klienta
	TerminalName string `soap:"terminal_name"`  // nazwa terminala
	ExtBillingId uint32 `soap:"ext_billing_id"` // id klienta w syst. lokalnym
	Caller       string `soap:"caller"`         // numer, z którego dzwoniono (format krajowy)
	BillCPB      string `soap:"bill_cpb"`       // numer, na który dzwoniono (format międzynarodowy bez +)
	CallType     string `soap:"calltype"`       // {incoming, outgoing, disa, forwarded, internal, vpbx}

	Country     string `soap:"country"`     // kierunek - kraj
	Description string `soap:"description"` // opis kierunku
	Operator    string `soap:"operator"`    // operator kierunku
	Type        string `soap:"type"`        // {mobile, geographic, premium, aus, NGN}

	Price      float64 `soap:"price"`       // cena połączenia dla klienta
	Cost       float64 `soap:"cost"`        // koszt połączenia dla klienta
	InitCharge float64 `soap:"init_charge"` // opłata za inicjacje połączenia dla klienta

	RePrice          float64 `soap:"reseller_price"`       // cena połączenia dla reselera
	ReCost           float64 `soap:"reseller_cost"`        // koszt połączenia dla reselera
	ReInitCharge     float64 `soap:"reseller_init_charge"` // opłata za inicjację połączenia dla reselera
	Margin           float64 `soap:"margin"`               // marża reselera
	SubscriptionUsed bool    `soap:"subscription_used"`    // czy został użyty abonament (np. darmowe minuty)

	PlatformType string `soap:"platform_type"` // typ usługi platformy Hiperus C5 {VOIP, WLR}
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

func (b *Billing) Scan(c *Call) error {
	row := b.rs.Children[b.n-1]
	err := row.LoadStruct(c, false)
	c.Duration *= time.Second
	return err
}

type getBilling struct {
	From         time.Time `soap:"from"`
	To           time.Time `soap:"to"`
	Offset       int       `soap:"offset"`
	Limit        *int      `soap:"limit"`
	Compress     bool      `soap:"compress"`
	SuccessCalls bool      `soap:"success_calls"`
	CallType     *string   `soap:"calltype"`
	CustomerId   *int      `soap:"id_customer"`
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
