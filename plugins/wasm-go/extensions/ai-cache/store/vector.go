package store

import (
	"fmt"
	"strings"
)

type Vector struct {
	Content []float32 `json:"content"`
	Length  uint32    `json:"length"`
	Raw     string    `json:"raw"`
	Answer  string    `json:"answer"`
}

func VectorToString(vector Vector) string {
	strSlice := make([]string, len(vector.Content))
	for i, v := range vector.Content {
		strSlice[i] = fmt.Sprintf("%f", v)
	}
	return fmt.Sprintf("[%s]", strings.Join(strSlice, ","))
}

func StringToVector(str string) Vector {
	str = strings.Trim(str, "[]")
	strs := strings.Split(str, ",")
	content := make([]float32, len(strs))
	for i, s := range strs {
		fmt.Sscanf(s, "%f", &content[i])
	}
	return Vector{Content: content}
}
