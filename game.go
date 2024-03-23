package main

func runGame() {
	channel := make(chan string)

	go ping(channel)
	go pong(channel)

	channel <- "ping"
}

func ping(channel chan string) {
	for {
		msg := <-channel
		println("Player 1:", msg)
		channel <- "pong"
	}
}

func pong(channel chan string) {
	for {
		msg := <-channel
		println("Player 2:", msg)
		channel <- "ping"
	}
}
