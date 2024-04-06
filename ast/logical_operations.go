package ast

func NOTOperation(a bool) bool {
	return !a
}

func ANDOperation(a bool, b bool) bool {
	return a && b
}

func NANDOperation(a bool, b bool) bool {
	return !ANDOperation(a, b)
}

func OROperation(a bool, b bool) bool {
	return a || b
}

func NOROperation(a bool, b bool) bool {
	return !OROperation(a, b)
}

func XOROperation(a bool, b bool) bool {
	return OROperation(a, b) && NANDOperation(a, b)
}

func XNOROperation(a bool, b bool) bool {
	return !XOROperation(a, b)
}
