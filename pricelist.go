package hiperus

import (
	"github.com/ziutek/soap"
)

type CustomerPricelist struct {
	Id                 uint32 `soap:"id"`
	Name               string `soap:"name"`
	ChargeInternalCall bool   `soap:"charge_internal_call"`
}

type CustomerPricelistList struct {
	rs *soap.Element
	n  int
}

func (pl *CustomerPricelistList) Next() bool {
	if pl.n == len(pl.rs.Children) {
		return false
	}
	pl.n++
	return true
}

func (pl *CustomerPricelistList) Scan(p *CustomerPricelist) error {
	row := pl.rs.Children[pl.n-1]
	return row.LoadStruct(p, false)
}

// GetCustomerPricelistList zwraca liste zdefiniowanych cenników
func (s *Session) GetCustomerPricelistList() (*CustomerPricelistList, error) {
	rs, err := s.cmd("GetCustomerPricelistList", struct{}{})
	if err != nil {
		return nil, err
	}
	return &CustomerPricelistList{rs, 0}, nil
}

// GetCustomerPricelis zwraca cennik o podanym id (jeśli różne od 0) lub nazwie
// (jeśli różna od ""). Zwraca cennik o id == 0 jeśli nic nie znalazła.
func (s *Session) GetCustomerPricelist(id uint32, name string) (
	p CustomerPricelist, err error) {

	pl, err := s.GetCustomerPricelistList()
	if err != nil {
		return
	}
	for pl.Next() {
		err = pl.Scan(&p)
		if err != nil {
			p = CustomerPricelist{}
			return
		}
		if id != 0 && p.Id == id {
			return
		}
		if name != "" && p.Name == name {
			return
		}
	}
	p = CustomerPricelist{}
	return
}
