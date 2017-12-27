// main.go
package main

import (
  "fmt"
  "net/http"

  "github.com/auth0/go-jwt-middleware"
  "github.com/codegangsta/negroni"
  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/mux"
)

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  token := r.Context().Value("user")

  /*
  type Token struct {
  Raw       string                 // The raw token.  Populated when you Parse a token
  Method    SigningMethod          // The signing method used or to be used
  Header    map[string]interface{} // The first segment of the token
  Claims    Claims                 // The second segment of the token
  Signature string                 // The third segment of the token.  Populated when you Parse a token
  Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}
  */
  
  jwttoken := token.(*jwt.Token)
  fmt.Println("token Header: ", jwttoken.Header)
  fmt.Println("token Claims: ", jwttoken.Claims)

  claims := jwttoken.Claims.(jwt.MapClaims)
  if val, ok := claims["username"]; ok{
    fmt.Println("username: ", val)
  } else {
    fmt.Println("username is not in jwt token payload")
    http.Error(w, "username not found in JWT token payload", http.StatusUnauthorized)
    return
  }

  fmt.Fprintf(w, "This is an authenticated request")
  fmt.Fprintf(w, "Claim content:\n")
  for k, v := range claims {
    fmt.Fprintf(w, "%s :\t%#v\n", k, v)
  }
})

func main() {
  r := mux.NewRouter()

  jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      return []byte("secret"), nil
    },
    // When set, the middleware verifies that tokens are signed with the specific signing algorithm
    // If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
    // Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
    SigningMethod: jwt.SigningMethodHS256,
  })

  r.Handle("/ping", negroni.New(
    negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
    negroni.Wrap(myHandler),
  ))
  http.Handle("/", r)
  http.ListenAndServe(":3001", nil)
}