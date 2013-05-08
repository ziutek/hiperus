package hiperus

import ()

type getFirstNumber struct {
	SN *string `soap:"sn"`
}

// GetFirstFreePlatformNumber zwraca pierwszy wolny numer w podanej strefie
// numeracyjnej sn. Jeśli nie ma wolnego numeru w podanej sn zwraca number == ""
func (s *Session) GetFirstFreePlatformNumber(sn string) (number, countryCode string, err error) {
	var arg getFirstNumber
	if sn != "" {
		arg.SN = &sn
	}
	rs, err := s.cmd("GetFirstFreePlatformNumber", arg)
	if err != nil {
		return
	}
	row, err := firstRow(rs)
	if err != nil {
		if err == ErrEmptyResultSet {
			err = nil
		}
		return
	}

	e, err := row.Get("free_number")
	if err != nil {
		return
	}
	number = e.AsStr()

	e, err = row.Get("country_code")
	if err != nil {
		return
	}
	countryCode = e.AsStr()
	return
}

type NumberData struct {
	Number      string `soap:"number"`
	CountryCode string `soap:"country_code"`
	SN          string `soap:"sn"`
	Secondary   bool   `soap:"is_main"` // false dla numeru głównego
	CLIR        bool   `soap:"clir"`
	VirtualFax  bool   `soap:"virtual_fax"`
}
