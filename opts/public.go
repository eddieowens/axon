package opts

type Opt[V any] func(opts *V)

func ApplyOpts[V any](defaults *V, opts ...Opt[V]) V {
	if defaults == nil {
		defaults = new(V)
	}

	for _, v := range opts {
		v(defaults)
	}

	return *defaults
}
