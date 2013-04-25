package hiperus

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/ziutek/soap"
	"io"
)

type soapEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	XSD string `xml:"xmlns:xsd,attr"`
	XSI string `xml:"xmlns:xsi,attr"`
	Enc string `xml:"xmlns:SOAP-ENC,attr"`

	EncStyle string `xml:"http://schemas.xmlsoap.org/soap/envelope/ encodingStyle,attr"`
	Body     soapBody
}

type soapBody struct {
	Fault    *soap.Fault
	Request  *soapRequest
	Response *soapResponse
}

type soapRequest struct {
	XMLName xml.Name `xml:"request"`

	Param0 *soap.Element
	Param1 *soap.Element
	Param2 *soap.Element
	Param3 *soap.Element
}

func newSoapRequest(p0, p1, p2, p3 interface{}) (*bytes.Buffer, error) {
	envelope := soapEnvelope{
		XSD:      "http://www.w3.org/2001/XMLSchema",
		XSI:      "http://www.w3.org/2001/XMLSchema-instance",
		Enc:      "http://schemas.xmlsoap.org/soap/encoding/",
		EncStyle: "http://schemas.xmlsoap.org/soap/encoding/",

		Body: soapBody{
			Request: &soapRequest{
				Param0: soap.MakeElement("param0", p0),
				Param1: soap.MakeElement("param1", p1),
				Param2: soap.MakeElement("param2", p2),
				Param3: soap.MakeElement("param3", p3),
			},
		},
	}
	out := bytes.NewBufferString(xml.Header)
	e := xml.NewEncoder(out)
	e.Indent("", "    ")
	err := e.Encode(&envelope)
	if err != nil {
		return nil, err
	}
	out.WriteByte('\n')
	//os.Stdout.Write(out.Bytes())
	return out, nil
}

type soapResponse struct {
	XMLName xml.Name `xml:"requestResponse"`

	Return *soap.Element `xml:"return"`
}

func parseResponse(r io.Reader) (*soap.Element, error) {
	envelope := new(soapEnvelope)
	if err := xml.NewDecoder(r).Decode(envelope); err != nil {
		return nil, err
	}
	body := envelope.Body
	if body.Fault != nil {
		return nil, body.Fault
	}
	if body.Response == nil {
		return nil, fmt.Errorf(
			"hiperus: SOAP envelope body doesn't contain response field",
		)
	}
	ret := body.Response.Return
	// Sprawdzamy czy operacja zakonczyla sie sukcesem
	success, err := ret.Get("success")
	if err != nil {
		return nil, err
	}
	ok, err := success.Bool()
	if err != nil {
		return nil, err
	}
	if !ok {
		errMsg, err := ret.Get("error_message")
		e, err := errMsg.Str()
		if err == nil {
			err = errors.New(e)
		}
		return nil, err
	}
	return ret, nil
}
