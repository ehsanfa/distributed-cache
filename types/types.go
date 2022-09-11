package types

type Req struct {
	Action int8
	Key string
	Value string
}

type Resp struct {
	Ok bool
	Key string
	Value string
}