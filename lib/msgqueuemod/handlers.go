package msgqueuemod

import "gitlab.com/snap-clickstaff/torque-go/lib/utils"

var MessageHandlerMap = make(map[MessageType]MessageHandler)

func RegisterHandler(msgType MessageType, handler MessageHandler) error {
	if _, ok := MessageHandlerMap[msgType]; ok {
		return utils.IssueErrorf("register a duplicated message queue handler key `%v`", msgType)
	}

	MessageHandlerMap[msgType] = handler
	return nil
}
