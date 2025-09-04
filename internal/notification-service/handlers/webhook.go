package handlers

import (
	"bytes"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/entity"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/senders"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/service"
	"io"
	"log"
	"net/http"
	"strings"
)

type WebhookHandler struct {
	sender         *senders.TelegramSender
	bindingService *service.BindingService
}

func NewWebhookHandler(sender *senders.TelegramSender, bindingService *service.BindingService) *WebhookHandler {
	return &WebhookHandler{
		sender:         sender,
		bindingService: bindingService,
	}
}
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	log.Printf("Incoming request body: %s", string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	switch r.Method {
	case http.MethodPost:
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			errorJSON(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		var update entity.TelegramUpdate
		err := decodeJSON(w, r, &update, 1<<20)
		if err != nil {
			respondDecodeError(w, err)
			return
		}

		log.Printf("Получено сообщение от %v: %q", update.Message.Chat.ID, update.Message.Text)

		if update.Message == nil || update.Message.Chat == nil {
			log.Println("Сообщение или чат пусты — не обрабатываем")
			return
		}

		chatId := update.Message.Chat.ID
		msg := update.Message.Text

		switch strings.ToLower(strings.TrimSpace(msg)) {
		case "/start":
			err := h.sender.SendMessage(chatId, "👋 Привет! Чтобы привязать Telegram к Task Manager — просто отправь свою почту.")
			if err != nil {
				log.Printf("ошибка отправки в Telegram: %v", err)
			}
			return
		case "/help":
			err := h.sender.SendMessage(chatId, "📖 Доступные команды:\n/start — инструкция\n/help — список команд\n📧 Также вы можете отправить свою почту для привязки.")
			if err != nil {
				log.Printf("ошибка отправки в Telegram: %v", err)
			}
			return
		}

		err = h.bindingService.BindEmailToChat(r.Context(), msg, chatId)
		if err != nil {
			log.Printf("ошибка привязки: %v", err)
			h.sender.SendMessage(chatId, "Ошибка при привязке почты. Попробуйте позже.")
			return
		}

		h.sender.SendMessage(chatId, "✅ Telegram успешно привязан к вашей почте!")

	}
}
