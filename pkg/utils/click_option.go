package utils

func ClickOption(options []string, separator string, cx int, offset int) (int, string) {
	sep := len(separator)
	sections := make([]int, 0)
	preFix := 0

	cx = cx - offset
	if cx < 0 {
		return -1, ""
	}

	var selected string
	for i, opt := range options {
		left := preFix + i*sep

		words := len([]rune(opt))

		right := left + words - 1
		preFix += words

		sections = append(sections, left, right)
	}

	optionIndex := -1
	for i := 0; i < len(sections); i += 2 {
		left := sections[i]
		right := sections[i+1]
		if cx >= left && cx <= right {
			optionIndex = i / 2
			selected = options[optionIndex]
			break
		}
	}
	return optionIndex, selected
}
