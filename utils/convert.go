package utils

import "github.com/jinzhu/copier"

func TransformDataOrPanic[T any](from any, to T, ignoreEmpty ...bool) T {
	isIgnoreEmpty := false

	if len(ignoreEmpty) > 0 {
		isIgnoreEmpty = ignoreEmpty[0]
	}

	err := copier.CopyWithOption(&to, from, copier.Option{
		DeepCopy:    false,
		IgnoreEmpty: isIgnoreEmpty,
	})
	LogAndPanicIfError(err, "failed when transform data")

	return to
}
