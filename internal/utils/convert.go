package utils

// BoolToString конвертирует булево значение в строку "1" или "0" для передачи в API РЖД. Почему они не используют булево значение в запросах?
func BoolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func BoolToYesNoLowerCase(b bool) string {
	if b {
		return "y"
	}
	return "n"
}
