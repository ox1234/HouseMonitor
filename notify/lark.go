package notify

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core"
	im "github.com/larksuite/oapi-sdk-go/service/im/v1"
	"houseMonitor/log"
)

// LarkNotifier ...
type LarkNotifier struct {
	ctx       context.Context
	chatID    string
	chatName  string
	imService *im.Service
}

// NewLarkNotify ...
func NewLarkNotify(appID string, appSecret string, chatName string) *LarkNotifier {
	appSettings := core.NewInternalAppSettings(
		core.SetAppCredentials(appID, appSecret))

	conf := core.NewConfig(core.DomainFeiShu, appSettings, core.SetLoggerLevel(core.LoggerLevelError))

	lark := &LarkNotifier{
		ctx:       context.Background(),
		chatName:  chatName,
		imService: im.NewService(conf),
	}

	chatID, err := lark.getChatID()
	if err != nil {
		log.Warn("get %s chat id fail: %s", chatName, err)
	}
	lark.chatID = chatID
	return lark
}

func (l *LarkNotifier) SendMessage(content string) error {
	coreCtx := core.WrapContext(l.ctx)
	reqCall := l.imService.Messages.Create(coreCtx, &im.MessageCreateReqBody{
		ReceiveId: l.chatID,
		Content:   content,
		MsgType:   "interactive",
	})
	reqCall.SetReceiveIdType("chat_id")
	_, err := reqCall.Do()
	if err != nil {
		return fmt.Errorf("lark send message fail: %w", err)
	}
	return nil
}

func (l *LarkNotifier) getChatID() (string, error) {
	coreCtx := core.WrapContext(l.ctx)
	reqCall := l.imService.Chats.List(coreCtx)
	msg, err := reqCall.Do()
	if err != nil {
		return "", fmt.Errorf("get chat id fail: %s", err)
	}

	for _, item := range msg.Items {
		if item.Name == l.chatName {
			return item.ChatId, nil
		}
	}

	return "", fmt.Errorf("no such chat %s", l.chatName)
}
