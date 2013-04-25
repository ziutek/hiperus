package hiperus

import (
	"github.com/ziutek/soap"
)

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
	*soap.Element, error) {

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
	return rs, nil
}
