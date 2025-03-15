package processor

// GetColumnLetter возвращает букву столбца по его номеру (1 -> A, 2 -> B, ...)
func GetColumnLetter(col int) string {
	return string(rune('A' + col - 1))
}
