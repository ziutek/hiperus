package hiperus

import (
	"github.com/ziutek/soap"
	"time"
)

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

type PSTNNumber struct {
	Number       string    `soap:"number"` //
	CountryCode  string    `soap:"country_code"`
	Extension    string    `soap:"extension"`
	IsMain       bool      `soap:"is_main"`
	CLIR         bool      `soap:"clir"`
	VirtualFax   bool      `soap:"virtual_fax"`
	TerminalId   uint32    `soap:"id"`            // identyfikator terminala SIP
	TerminalName string    `soap:"terminal_name"` // nazwa terminala SIP
	AuthId       uint32    `soap:"id_auth"`
	TempNumber   bool      `soap:"temp_number"`
	DISA         bool      `soap:"disa_enabled"`
	Voicemail    bool      `soap:"voicemail_enabled"`
	CreateDate   time.Time `soap:"create_date"`
}

type PSTNNumberList struct {
	rs *soap.Element
	n  int
}

func (nl *PSTNNumberList) Next() bool {
	if nl.n == len(nl.rs.Children) {
		return false
	}
	nl.n++
	return true
}

func (nl *PSTNNumberList) Scan(pn *PSTNNumber) (err error) {
	row := nl.rs.Children[nl.n-1]
	return row.LoadStruct(pn, false)
}

type getPSTNNumberList struct {
	CustomerId uint32 `soap:"id_customer"`
	Offset     int    `soap:"offset"`
	Limit      *int   `soap:"limit"`
}

func (s *Session) GetPSTNNumberList(customerId uint32, offset, limit int) (
	*PSTNNumberList, error) {

	arg := getPSTNNumberList{
		CustomerId: customerId,
		Offset:     offset,
	}
	if limit != 0 {
		arg.Limit = &limit
	}
	rs, err := s.cmd("GetExtensionList", arg)
	if err != nil {
		return nil, err
	}
	return &PSTNNumberList{rs, 0}, nil
}

type addExtension struct {
	CustomerId  uint32 `soap:"id_customer"`
	TerminalId  uint32 `soap:"id_terminal"`
	Number      string `soap:"number"`
	CountryCode string `soap:"country_code"`
	IsMain      bool   `soap:"is_main"`
	CLIR        bool   `soap:"clir"`
	VirtualFax  bool   `soap:"virtual_fax"`
}

func (s *Session) AddExtension(customerId, terminalId uint32, number,
	countryCode string, isMain, clir, virtualFax bool) (uint32, error) {

	arg := addExtension{
		CustomerId:  customerId,
		TerminalId:  terminalId,
		Number:      number,
		CountryCode: countryCode,
		IsMain:      isMain,
		CLIR:        clir,
		VirtualFax:  virtualFax,
	}

	rs, err := s.cmd("AddExtension", &arg)
	if err != nil {
		return 0, err
	}
	e, err := firstRow(rs)
	if err != nil {
		return 0, err
	}
	if e, err = e.Get("id_extension"); err != nil {
		return 0, err
	}
	return e.AsUint32()
}
