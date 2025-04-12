package whatsapp

import (
	"database/sql"
	"log/slog"
	"reflect"

	"github.com/weilence/whatsapp-client/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type eventHandler struct {
	client  *whatsmeow.Client
	db      *sql.DB
	handler func(v any)
}

func newEventHandler(client *whatsmeow.Client, db *sql.DB) *eventHandler {
	return &eventHandler{
		client: client,
		db:     db,
	}
}

func (h *eventHandler) handleEvent(v any) {
	if *config.Env == "dev" {
		slog.Info("handle event", "type", reflect.TypeOf(v))
	}

	switch evt := v.(type) {
	case *events.HistorySync:
		h.handleHistorySync(evt)
	case *events.Message:
		h.handleMessage(evt)
	}
}

func (h *eventHandler) handleMessage(evt *events.Message) {
	if evt.Info.Type == "reaction" {
		return
	}

	if evt.Info.Edit == types.EditAttributeSenderRevoke || evt.Info.Edit == types.EditAttributeAdminRevoke {
		err := revokeMessage(h.db, evt.Info.ID)
		if err != nil {
			slog.Error("revoke message", "err", err)
		}
		return
	}

	err := insertMessage(h.db, evt.Info.ID, evt.Info.Sender.String(), evt.Info.Chat.String(), evt.Message.GetConversation())
	if err != nil {
		slog.Error("insert message", "err", err)
	}
}

func (h *eventHandler) handleHistorySync(evt *events.HistorySync) {
	for _, conv := range evt.Data.GetConversations() {
		for _, historymsg := range conv.GetMessages() {
			chatJID, _ := types.ParseJID(conv.GetID())
			evt, err := h.client.ParseWebMessage(chatJID, historymsg.GetMessage())
			if err != nil {
				slog.Error("parse web message", "err", err)
				continue
			}

			h.handleMessage(evt)
		}
	}
}
