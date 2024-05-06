package utils

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func ConvertInts[U, T Int](s []T) (out []U) {
	out = make([]U, len(s))
	for i := range s {
		out[i] = U(s[i])
	}
	return out
}
