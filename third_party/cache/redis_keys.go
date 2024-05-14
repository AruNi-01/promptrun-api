package cache

const (
	Split = ":"
)

const (
	PrefixTicket = "ticket"
	PrefixPrompt = "prompt"
)

func TicketKey(ticket string) string {
	return PrefixTicket + Split + ticket
}

func PromptScoreChangeKey() string {
	return PrefixPrompt + Split + "score_change"
}
