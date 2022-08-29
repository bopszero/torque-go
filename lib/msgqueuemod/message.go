package msgqueuemod

import (
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

func NewMessage(msgType MessageType, data interface{}) Message {
	return Message{
		ID:   comutils.NewUUID4().String(),
		Type: msgType,
		Data: data,
	}
}

func NewMessageJSON(msgType MessageType, data interface{}) (Message, error) {
	dataJSON, err := comutils.JsonEncode(data)
	if err != nil {
		return Message{}, err
	}

	return NewMessage(msgType, dataJSON), nil
}

func NewMessageJsonF(msgType MessageType, data interface{}) Message {
	message, err := NewMessageJSON(msgType, data)
	comutils.PanicOnError(err)

	return message
}
