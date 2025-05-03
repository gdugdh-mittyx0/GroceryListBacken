package utils

import "context"

func LogWithUser(ctx context.Context, msg string) string {
	// session, ok := ctx.Value(middleware.SessionContextKey).(JWTClaims)
	// if ok {
	// 	msg = fmt.Sprintf("%s | role: %s | id_user: %d", msg, session.UserType, session.UserID)
	// }
	return msg
}
