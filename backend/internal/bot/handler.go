package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ai-bot-expenses/backend/internal/llm"
	"github.com/ai-bot-expenses/backend/internal/models"
	"github.com/ai-bot-expenses/backend/internal/service"
)

// Handler processes messages received by the Telegram bot.
type Handler struct {
	service   *service.ExpenseService
	llmClient *llm.Client
}

// NewHandler creates a new bot message handler.
func NewHandler(svc *service.ExpenseService, llmClient *llm.Client) *Handler {
	return &Handler{
		service:   svc,
		llmClient: llmClient,
	}
}

// HandleMessage parses and executes user commands from natural language text.
func (h *Handler) HandleMessage(ctx context.Context, userID int64, text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return "Пожалуйста, напишите сообщение о расходах."
	}

	if text == "/start" || text == "/help" {
		return "👋 Привет! Я бот для учёта расходов.\n\n" +
			"Просто пишите мне на обычном языке, например:\n" +
			"• «кофе 350» или «такси 1200 вчера» — добавить расход\n" +
			"• «сколько потратил на еду в этом месяце» — узнать сумму\n" +
			"• «покажи траты за неделю» — список последних трат\n" +
			"• «последнее было не 350, а 450» — изменить трату\n" +
			"• «удали такси» — удалить трату\n\n" +
			"Вам доступны только ваши личные расходы!"
	}

	// Request LLM to understand intent
	intent, err := h.llmClient.ParseIntent(ctx, text)
	if err != nil {
		return fmt.Sprintf("⚠️ Извините, не удалось распознать сообщение: %v", err)
	}

	switch intent.Action {
	case llm.ActionCreate:
		return h.handleCreate(userID, intent)
	case llm.ActionList:
		return h.handleList(userID, intent)
	case llm.ActionUpdate:
		return h.handleUpdate(userID, intent)
	case llm.ActionDelete:
		return h.handleDelete(userID, intent)
	default:
		if intent.Message != "" {
			return intent.Message
		}
		return "Я не понял ваше сообщение. Попробуйте написать, например: «кофе 350» или «сколько потратил на еду в этом месяце»."
	}
}

func (h *Handler) handleCreate(userID int64, intent *llm.Response) string {
	if intent.Amount == nil || *intent.Amount <= 0 {
		return "⚠️ Не удалось определить сумму расхода. Пожалуйста, укажите число, например: «обед 500»."
	}

	category := string(models.CategoryOther)
	if intent.Category != nil && models.IsValidCategory(*intent.Category) {
		category = *intent.Category
	}

	spentAt := time.Now().Format("2006-01-02")
	if intent.SpentAt != nil && *intent.SpentAt != "" {
		if _, err := time.Parse("2006-01-02", *intent.SpentAt); err == nil {
			spentAt = *intent.SpentAt
		}
	}

	comment := ""
	if intent.Comment != nil {
		comment = strings.TrimSpace(*intent.Comment)
	}
	if comment == "" {
		comment = category
	}

	req := models.CreateExpenseRequest{
		UserID:   userID,
		Amount:   *intent.Amount,
		Category: category,
		SpentAt:  spentAt,
		Comment:  comment,
	}

	expense, err := h.service.Create(req)
	if err != nil {
		return fmt.Sprintf("❌ Ошибка при сохранении траты: %v", err)
	}

	return fmt.Sprintf("✅ Записал: %s — %.2f ₽\nКатегория: %s\nДата: %s",
		expense.Comment, expense.Amount, expense.Category, expense.SpentAt.Format("02.01.2006"))
}

func (h *Handler) handleList(userID int64, intent *llm.Response) string {
	filter := models.ExpenseFilter{
		UserID: &userID,
	}

	if intent.Category != nil && *intent.Category != "" && models.IsValidCategory(*intent.Category) {
		filter.Category = intent.Category
	}

	if intent.Period != nil {
		if intent.Period.From != "" {
			filter.DateFrom = &intent.Period.From
		}
		if intent.Period.To != "" {
			filter.DateTo = &intent.Period.To
		}
	}

	expenses, total, err := h.service.List(filter)
	if err != nil {
		return fmt.Sprintf("❌ Ошибка при получении трат: %v", err)
	}

	if total == 0 {
		return "🔍 Трат по вашему запросу не найдено."
	}

	totalSum, err := h.service.SumByFilter(userID, filter)
	if err != nil {
		totalSum = 0
		for _, e := range expenses {
			totalSum += e.Amount
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 Найдено трат: %d (Всего: %.2f ₽)\n\n", total, totalSum))

	maxItems := 10
	for i, e := range expenses {
		if i >= maxItems {
			sb.WriteString(fmt.Sprintf("\n... и ещё %d трат", int(total)-maxItems))
			break
		}
		sb.WriteString(fmt.Sprintf("• %s: %s — %.2f ₽ (%s)\n",
			e.SpentAt.Format("02.01"), e.Comment, e.Amount, e.Category))
	}

	return sb.String()
}

func (h *Handler) handleUpdate(userID int64, intent *llm.Response) string {
	var targetExpense *models.Expense
	var err error

	target := ""
	if intent.Target != nil {
		target = strings.ToLower(strings.TrimSpace(*intent.Target))
	}

	if target == "last" || target == "" {
		targetExpense, err = h.service.GetLastByUserID(userID)
	} else {
		targetExpense, err = h.service.FindByUserIDAndComment(userID, target)
	}

	if err != nil || targetExpense == nil {
		return "⚠️ Не удалось найти подходящую трату для изменения. Проверьте правильность запроса."
	}

	updateReq := models.UpdateExpenseRequest{}
	if intent.Amount != nil && *intent.Amount > 0 {
		updateReq.Amount = intent.Amount
	}
	if intent.Category != nil && models.IsValidCategory(*intent.Category) {
		updateReq.Category = intent.Category
	}
	if intent.SpentAt != nil && *intent.SpentAt != "" {
		if _, err := time.Parse("2006-01-02", *intent.SpentAt); err == nil {
			updateReq.SpentAt = intent.SpentAt
		}
	}
	if intent.Comment != nil && *intent.Comment != "" {
		updateReq.Comment = intent.Comment
	}

	updated, err := h.service.Update(targetExpense.ID, updateReq)
	if err != nil {
		return fmt.Sprintf("❌ Ошибка при обновлении траты: %v", err)
	}

	return fmt.Sprintf("✏️ Трата обновлена:\n%s — %.2f ₽\nКатегория: %s\nДата: %s",
		updated.Comment, updated.Amount, updated.Category, updated.SpentAt.Format("02.01.2006"))
}

func (h *Handler) handleDelete(userID int64, intent *llm.Response) string {
	var targetExpense *models.Expense
	var err error

	target := ""
	if intent.Target != nil {
		target = strings.ToLower(strings.TrimSpace(*intent.Target))
	}

	if target == "last" || target == "" {
		targetExpense, err = h.service.GetLastByUserID(userID)
	} else {
		targetExpense, err = h.service.FindByUserIDAndComment(userID, target)
	}

	if err != nil || targetExpense == nil {
		return "⚠️ Не удалось найти трату для удаления."
	}

	deletedComment := targetExpense.Comment
	deletedAmount := targetExpense.Amount
	deletedDate := targetExpense.SpentAt.Format("02.01.2006")

	if err := h.service.Delete(targetExpense.ID); err != nil {
		return fmt.Sprintf("❌ Ошибка при удалении траты: %v", err)
	}

	return fmt.Sprintf("🗑 Трата успешно удалена:\n%s — %.2f ₽ от %s",
		deletedComment, deletedAmount, deletedDate)
}
