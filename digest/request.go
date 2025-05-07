package digest

// request is private, just works as a internal interface
type request[T any] interface {
	GetRequest() (T, error)
}

type RequestAuth interface {
	request[*Authentication]
	SendRequest() (string, error)
}

type RequestDownload interface {
}

type RequestRequest interface {
}

type RequestVerify interface {
}
