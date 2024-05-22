package hamming

type Usecase interface {
	ChannelTransmit(data string) (string, bool, bool, error)
}
