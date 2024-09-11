package begger

func ValueToPointer[T any](value T) *T { return &value }
