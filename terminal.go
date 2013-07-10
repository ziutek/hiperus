package hiperus

import (
	"github.com/ziutek/soap"
	"time"
)

type Terminal struct {
	Id uint32 `soap:"id,in"` // identifikator terminala

	Name          string `soap:"username"`       // nazwa terminala
	Password      string `soap:"password"`       // hasło terminala
	ScreenNumbers bool   `soap:"screen_numbers"` // czy sprawdzać przesyłne num.
	T38Fax        bool   `soap:"t38_fax"`        // czy obsługiwać T38 fax

	CustomerId     uint32 `soap:"id_customer"`
	PriceListId    uint32 `soap:"id_pricelist"`
	SubscriptionId uint32 `soap:"id_subscription,omitempty"` // id abonamentu

	CustomerName  string `soap:"customer_name,in"`
	PriceListName string `soap:"pricelist_name,in"`

	BalanceValue float64 `soap:"balance_value,in"` // środki na koncie

	AuthId           uint32    `soap:"id_auth,in"`
	SubscriptionFrom time.Time `soap:"subscription_from,in"` // od kiedy abon.
	SubscriptionTo   time.Time `soap:"subscription_to,in"`   // do kiedy abon.
	ValueLeft        float64   `soap:"value_left","`         // pozost. w abon.

	LocationId uint32 `soap:"id_terminal_location,in"` // id lokalizacji term.
	AreaCode   string `soap:"area_code,in"`            // strefa numeracyjna
	Borough    string `soap:"borough,in"`              // gmina dla poł. alarm.
	County     string `soap:"county,in"`               // powiad dla poł. alarm.
	Province   string `soap:"province,in"`             // województwo

	SIPProxy string `soap:"sip_proxy,in"`
}

func (s *Session) AddTerminal(t *Terminal) (uint32, error) {
	rs, err := s.cmd("AddTerminal", t)
	if err != nil {
		return 0, err
	}
	e, err := firstRow(rs)
	if err != nil {
		return 0, err
	}
	if e, err = e.Get("id_terminal"); err != nil {
		return 0, err
	}
	return e.AsUint32()
}

type terminalId struct {
	TerminalId uint32 `soap:"id_terminal"`
}

func (s *Session) DelTerminal(id uint32) error {
	_, err := s.cmd("DelTerminal", terminalId{id})
	return err
}

type TerminalList struct {
	rs *soap.Element
	n  int
}

func (tl *TerminalList) Next() bool {
	if tl.n == len(tl.rs.Children) {
		return false
	}
	tl.n++
	return true
}

func (tl *TerminalList) Scan(t *Terminal) (err error) {
	row := tl.rs.Children[tl.n-1]
	return row.LoadStruct(t, false)
}

type getTerminalList struct {
	CustomerId uint32 `soap:"id_customer"`
	Offset     int    `soap:"offset"`
	Limit      *int   `soap:"limit"`
}

func (s *Session) GetTerminalList(customerId uint32, offset, limit int) (
	*TerminalList, error) {

	arg := getTerminalList{
		CustomerId: customerId,
		Offset:     offset,
	}
	if limit != 0 {
		arg.Limit = &limit
	}
	rs, err := s.cmd("GetTerminalList", arg)
	if err != nil {
		return nil, err
	}
	return &TerminalList{rs, 0}, nil
}
