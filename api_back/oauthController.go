package api_back

import (
	"GoMeLi/api_front"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)


var code string // Esta variable contendra el codigo para obtener el accses token y posteriormente el refresh_token

var AccessToken string // se almacenara el access token para las posteriores transacciones

var url string = "http://localhost:80/webtest/oauth" // URL puesta en la app de MeLi

// struct que se enviara como body para obtener el Access token
type Token struct {
	GrantType string 		`json:"grant_type"`
	ClientId int 			`json:"client_id"`
	ClientSecret string 	`json:"client_secret"`
	Code string				`json:"code"`
	RedirectUri string 		`json:"redirect_uri,omitempty"` //el ",omitempty" sirve para que al pasar el struct a
															//JSON, si la var esta vacia se omite y no se envia
}

// struct para almacenar la respuesta de MeLi
type TokenResp struct {
	AccessToken string 		`json:"access_token"`
	TokenType string 		`json:"token_type"`
	ExpiresIn int 			`json:"expires_in"`
	Scope string 			`json:"scope"`
	UserId int 				`json:"user_id"`
	RefreshToken string 	`json:"refresh_token"`
}


func GetCode(c *gin.Context){
	// obtenemos el codigo de intercambio y nos aseguramos de que no este vacio
	code = c.Query("code")
	if code == "" {
		c.String(400, "HTTP 400 Missing param code")
		return
	}

	// mostramos la pagina de home
	api_front.HomePage(c)

	// llamamos a la funcion para obtener el token por primera vez
	tokenRequest(true)
}

func tokenRequest( firstChange bool ) {

	b, err := json.Marshal(bodyToken( firstChange ))

	// comprobamos que no haya un error en la conversion
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))

	// hacemos el post del body
	resp, err := http.Post("https://api.mercadolibre.com/oauth/token",
					 "application/json; application/x-www-form-urlencoded",
								bytes.NewBuffer(b))

	if err != nil {
		fmt.Errorf("Error ",err.Error())
		return
	}

	defer resp.Body.Close()

	// leemos la respuesta de MeLi
	data, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(data))

	// decodificamos la respuesta y la almacenamos en una tokenResp
	var tokenResp TokenResp
	json.Unmarshal(data, &tokenResp)
	fmt.Printf("%+v\n", tokenResp)

	//extraemos el access token
	AccessToken = tokenResp.AccessToken

	//extraemos el refresh token
	code=tokenResp.RefreshToken

	fmt.Printf(AccessToken)
}

func bodyToken( firstChange bool ) Token {

	if(firstChange) {
		return Token{	GrantType: "authorization_code",
						ClientId: 6719038448258240,
						ClientSecret: "qmxiwj6zMUkNyWs1YzdOHkuCkkquJfVw",
						Code: code,
						RedirectUri: url}
	}

	return Token{	GrantType: "refresh_code",
					ClientId: 6719038448258240,
					ClientSecret: "qmxiwj6zMUkNyWs1YzdOHkuCkkquJfVw",
					Code: code}

}





