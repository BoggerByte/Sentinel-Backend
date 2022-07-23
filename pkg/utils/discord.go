package utils

func AnyOfPermissions(p1 int64, p2 int64) bool {
	return p1&p2 > 0
}
