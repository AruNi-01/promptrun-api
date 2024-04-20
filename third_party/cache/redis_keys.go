package cache

const (
	Split = ":"
)

const (
	PrefixTicket = "ticket"
)

func TicketKey(ticket string) string {
	return PrefixTicket + Split + ticket
}
