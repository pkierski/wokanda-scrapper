package trialdownloader

import "iter"

func MutableValues[Slice ~[]E, E any](s Slice) iter.Seq[*E] {
	return func(yield func(*E) bool) {
		for i := range s {
			if !yield(&s[i]) {
				return
			}
		}
	}
}
