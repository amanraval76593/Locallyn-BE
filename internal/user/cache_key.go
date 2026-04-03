package user

import "fmt"

func verifyUserKey(code string) string {
	return fmt.Sprintf("verify-user:%s", code)
}
