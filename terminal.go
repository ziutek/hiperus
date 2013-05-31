package hiperus

import (
	"github.com/ziutek/soap"
	"time"
)

type Terminal struct {
	Id uint32 `soap:"id"` // identifikator terminala

	Name          string `soap:"username"`       // nazwa terminala
	Password      string `soap:"password"`       // hasło terminala
	ScreenNumbers bool   `soap:"screen_numbers"` // czy sprawdzać przesyłne num.
	T38Fax        bool   `soap:"t38_fax"`        // czy obsługiwać T38 fax

	CustomerId  uint32 `soap:"id_customer"`
	PriceListId uint32 `soap:"id_pricelist"`

	CustomerName  string `soap:"customer_name,in"`
	PriceListName string `soap:"pricelist_name,in"`

	BalanceValue     float64   `soap:"balance_value,in"`
	AuthId           uint32    `soap:"id_auth,in"`
	SubscriptionId   uint32    `soap:"id_subscription,in"`
	SubscriptionFrom time.Time `soap:"subscription_from,in"`
	SubscriptionTo   time.Time `soap:"subscription_to,in"`
	ValueLeft        float64   `soap:"value_left","`
	
	// TODO
}

/*func (s *Session) AddTerminal(t *Terminal) {
	rs, err := s.cmd("AddTerminal", t)
}*/

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
