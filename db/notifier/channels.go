package notifier

type Channel string

var (
	NewMessage    Channel = "new_message"
	DeleteMessage Channel = "delete_message"
)
