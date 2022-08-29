package emailmod

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gopkg.in/gomail.v2"
)

type Email struct {
	gomail.Message
}

func NewEmail() *Email {
	msg := gomail.NewMessage()
	return &Email{*msg}
}

func (this *Email) SetSubject(subject string) *Email {
	this.SetHeader(constants.EmailHeaderSubject, subject)
	return this
}

func (this *Email) SetFromAddress(address string) *Email {
	this.SetHeader(constants.EmailHeaderFrom, address)
	return this
}

func (this *Email) SetFromContact(name, address string) *Email {
	this.SetFromAddress(this.FormatAddress(address, name))
	return this
}

func (this *Email) AddToAddresses(addresses ...string) *Email {
	curToAddresses := this.GetHeader(constants.EmailHeaderTo)
	this.SetHeader(constants.EmailHeaderTo, append(curToAddresses, addresses...)...)
	return this
}

func (this *Email) AddToAddress(address string) *Email {
	this.AddToAddresses(address)
	return this
}

func (this *Email) AddToContact(name, address string) *Email {
	this.AddToAddress(this.FormatAddress(address, name))
	return this
}
