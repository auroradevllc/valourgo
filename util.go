package valourgo

func Ref[V any](v V) *V {
	return &v
}
