package runner

import (
	"context"
	"time"

	"google.golang.org/api/idtoken"
)

var (
	token          string
	expirationDate time.Time
)

func generateJWT(aud string) (string, error) {
	if token != "" && expirationDate.Sub(time.Now()).Seconds() > 100 {
		return token, nil
	}
	ctx := context.Background()
	tokenSource, err := idtoken.NewTokenSource(ctx, aud)
	if err != nil {
		return "", err
	}
	t, err := tokenSource.Token()
	if err != nil {
		return "", err
	}
	token = t.AccessToken
	expirationDate = time.Now().Add(time.Hour)
	return token, nil

}
