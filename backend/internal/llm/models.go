package llm

// Action defines the intent parsed from the user's message.
type Action string

const (
	ActionCreate  Action = "create"
	ActionList    Action = "list"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionUnknown Action = "unknown"
)

// Period represents a date range for queries.
type Period struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Response represents the structured output expected from the LLM.
type Response struct {
	Action   Action   `json:"action"`
	Amount   *float64 `json:"amount"`
	Category *string  `json:"category"`
	SpentAt  *string  `json:"spent_at"` // Format: YYYY-MM-DD
	Comment  *string  `json:"comment"`
	Period   *Period  `json:"period"`
	Target   *string  `json:"target"` // e.g. "last" or "такси"
	Message  string   `json:"message"`
}

// ChatMessage represents a single message in the chat completions payload.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the payload for an OpenAI-compatible /chat/completions endpoint.
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

// ChatChoice represents a single choice in the response.
type ChatChoice struct {
	Index   int         `json:"index"`
	Message ChatMessage `json:"message"`
}

// ChatResponse represents the response from an OpenAI-compatible /chat/completions endpoint.
type ChatResponse struct {
	ID      string       `json:"id"`
	Choices []ChatChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}
