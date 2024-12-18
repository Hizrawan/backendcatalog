package messaging

type Clients struct {
	SMS   *Every8dClient
	Email *MailgunClient
}
