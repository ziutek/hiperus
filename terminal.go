package hiperus

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

	BalanceValue float64 `soap:"balance_value,in"`
	AuthId       uint32  `soap:"id_auth"`

	//TODO
}

func (s *Session) AddTerminal(t *Terminal) {
	rs, err := s.cmd("AddTerminal", t)
}
