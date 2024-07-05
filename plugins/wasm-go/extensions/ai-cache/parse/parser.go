package parse

type Parser interface {
	// Convert 把用户输入的语句转换为多个完整语义的独立语句
	Convert(raw string) ([]string, error)
}
