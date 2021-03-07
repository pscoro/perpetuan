package main

type message struct {
	Sender player
	Text   string
}

func newMessage(sender player, text string) message {
	return message{
		Sender: sender,
		Text:   text,
	}
}
