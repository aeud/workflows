package runner

import (
	"context"
	"time"

	"google.golang.org/api/idtoken"
)

// const googlAuthURL = "https://www.googleapis.com/oauth2/v4/token"

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
	// log.Printf("Refreshing the authorization token")
	// ctx := context.Background()
	// creds, err := google.FindDefaultCredentials(ctx)
	// if err != nil {
	// 	return "", err
	// }
	// conf, err := google.JWTConfigFromJSON(creds.JSON)
	// if err != nil {
	// 	return "", err
	// }

	// cs := &jws.ClaimSet{
	// 	Iss: conf.Email,
	// 	Sub: conf.Email,
	// 	Aud: googlAuthURL,
	// }
	// hdr := &jws.Header{
	// 	Algorithm: "RS256",
	// 	Typ:       "JWT",
	// 	KeyID:     conf.PrivateKeyID,
	// }

	// privateKey := ParseKey(conf.PrivateKey)

	// jwToken, err := tokenGenerator(aud, cs, hdr, privateKey)
	// if err != nil {
	// 	return "", err
	// }
	// token = jwToken
	// log.Printf("Next expiration date: %s", expirationDate)
	// return jwToken, nil

}

// func ParseKey(key []byte) *rsa.PrivateKey {
// 	block, _ := pem.Decode(key)
// 	if block != nil {
// 		key = block.Bytes
// 	}
// 	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
// 	if err != nil {
// 		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
// 		if err != nil {
// 			return nil
// 		}
// 	}
// 	parsed, ok := parsedKey.(*rsa.PrivateKey)
// 	if !ok {
// 		return nil
// 	}
// 	return parsed
// }

// func tokenGenerator(aud string, cs *jws.ClaimSet, hdr *jws.Header, key *rsa.PrivateKey) (string, error) {
// 	iat := time.Now()
// 	exp := iat.Add(time.Hour)
// 	expirationDate = exp

// 	cs.PrivateClaims = map[string]interface{}{"target_audience": aud}
// 	cs.Iat = iat.Unix()
// 	cs.Exp = exp.Unix()

// 	msg, err := jws.Encode(hdr, cs, key)
// 	if err != nil {
// 		return "", err
// 	}

// 	f := url.Values{
// 		"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
// 		"assertion":  {msg},
// 	}

// 	res, err := http.PostForm(googlAuthURL, f)
// 	if err != nil {
// 		return "", err
// 	}
// 	c, err := ioutil.ReadAll(res.Body)
// 	defer res.Body.Close()
// 	if err != nil {
// 		return "", err
// 	}

// 	type resIDToken struct {
// 		IDToken string `json:"id_token"`
// 	}

// 	id := &resIDToken{}

// 	if err := json.Unmarshal(c, id); err != nil {
// 		return "", err
// 	}

// 	return id.IDToken, nil
// }
