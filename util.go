package valour

func Ref[V any](v V) *V {
	return &v
}
