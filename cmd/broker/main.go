package main
import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/api"
	"github.com/luideoz/API_Go/config"
	"github.com/luideoz/API_Go/models"
    "io"
)

// Estrucutra para la respuesta del endpoint de listar archivos
type FileResponse struct {
    Files []models.File `json:"files"`
}
func main() {
    config.Load("config/config.toml")
    api := gin.Default()
    api.GET("/version", version)
    api.POST("/signup", signup)
    api.GET("/login", login)

    api.POST("/:username/:doc_id", upload)
    api.GET("/:username/:doc_id", getFile)
    api.PUT("/:username/:doc_id", update)
    api.DELETE("/:username/:doc_id", delete)
    api.GET(("/:username/_all_docs"), listFiles)

    api.Run(":8080")
}

//version -> devuelve la version del broker
func version(c *gin.Context) {
    c.JSON(200, gin.H{
        "version": "1.0.0",
    })
}

//signup -> endpoint para registrar un usuario
func signup(c *gin.Context) {
   AuthRequest(c, "POST")
}

//login -> endpoint para loguear un usuario
func login(c *gin.Context) {
    AuthRequest(c, "GET")
}

//upload -> endpoint para subir un archivo
func upload(c *gin.Context) {
    FileRequest(c, "POST")
}

//getFile -> endpoint para obtener un archivo
func getFile(c * gin.Context) {
    FileRequest(c, "GET")
}

//update -> endpoint para actualizar un archivo
func update(c *gin.Context) {
    FileRequest(c, "PUT")
}

//delete -> endpoint para borrar un archivo
func delete(c *gin.Context) {
    FileRequest(c, "DELETE")
}

//listFiles -> endpoint para listar los archivos de un usuario
func listFiles(c *gin.Context) {
    FileRequest(c, "LIST")
}

//checkBody_user -> funcion para comprobar el cuerpo del mensaje en los endpoints de usuario
func checkBody_user(c *gin.Context) (bool,models.User,string) {
	//miramos si el cuerpo del mensaje esta vacio
	var user models.User
	if c.Request.Body == nil {
		return false,user,"empty body"
	}
	err := c.BindJSON(&user)
	if err != nil {
		return false,user,"invalid json"
	}
	//miramos si el campo username esta vacio
	if user.USERNAME == "" {
		return false,user,"empty username"
	}
	//miramos si el campo password esta vacio
	if user.PASSWORD == "" {
		return false,user,"empty password"
	}
	return true,user,""
}

//checkBody_file -> funcion para comprobar el cuerpo del mensaje en los endpoints de archivo
func checkBody_file(c *gin.Context) (bool,models.Json,string) {
	var json models.Json
	if c.Request.Body == nil {
		return false, json, "empty body"
	}
	err := c.BindJSON(&json)
	if err != nil {
		return false, json, "invalid json"
	}
	if json.Doc_content == "" {
		return false, json, "empty doc_content"
	}
	return true, json, ""
}

//checkHeader -> funcion para comprobar la cabecera del mensaje
func checkHeader(c *gin.Context) (bool,string) {
    token := c.Request.Header.Get("Authorization")
    if token == "" {
        return false, ""
    }
    return true, token
}

//AuthRequest -> funcion para realizar peticiones a los endpoints de autenticacion
func AuthRequest(c *gin.Context, method string) {
    valid, user, message := checkBody_user(c)
    if !valid {
        c.JSON(400, gin.H{"error": message})
        return
    }

    user.USERNAME = api.Encrypt_hash(user.USERNAME)
    user.PASSWORD = api.Encrypt_hash(user.PASSWORD)

    jsonData, err := json.Marshal(user)
    if err != nil {
        c.JSON(500, gin.H{"error": "internal error, json marshal"})
        return
    }
    if method == "POST" {
        response, err := http.Post("http://10.0.2.3:8081/signup", "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            c.JSON(500, gin.H{"error": "could not connect to auth service"})
            return
        }
        defer response.Body.Close()
        var data map[string]string
        err = json.NewDecoder(response.Body).Decode(&data)
        if err != nil {
            c.JSON(500, gin.H{"error": "internal error, json decode"})
            return
        }
        c.JSON(response.StatusCode, data)
        return
    }
    response, err := http.Post("http://10.0.2.3:8081/login", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        c.JSON(500, gin.H{"error": "could not connect to auth service"})
        return
    }
    defer response.Body.Close()
    var data map[string]string
    err = json.NewDecoder(response.Body).Decode(&data)
    if err != nil {
        c.JSON(500, gin.H{"error": "internal error, json decode"})
        return
    }
    c.JSON(response.StatusCode, data)

}

//FileRequest -> funcion para realizar peticiones a los endpoints de archivo
func FileRequest(c *gin.Context, method string)  { //username, doc_id, token

    //obtenecion de los parametros y encriptacion del nombre de usuario
    username := c.Param("username")
    doc_id := c.Param("doc_id")
    username = api.Encrypt_hash(username)

    //validacion de la cabecera
    valid, token := checkHeader(c)
    if !valid {
        c.JSON(400, gin.H{"error": "empty token"})
        return 
    }

    //peticion al servicio de autenticacion para la comprobacion del token
    response, err := http.Get("http://10.02.3:8081/"+username+"/"+token)
    if err != nil {
        c.JSON(500, gin.H{"error": "could not connect to auth service"})
        return
    }
    defer response.Body.Close()
    if response.StatusCode != 204 {
        body, _ := io.ReadAll(response.Body)
        c.JSON(response.StatusCode, gin.H{"error": string(body)})
        return
    }

    // peticion get al servicio de archivos
    if method == "GET" {
        response, err = http.Get("http://10.0.2.4:8082/"+username+"/"+doc_id)
        if err != nil {
            c.JSON(500, gin.H{"error": "could not connect to file service"})
            return
        }
        defer response.Body.Close()
        var jsonData models.Json
        err = json.NewDecoder(response.Body).Decode(&jsonData)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(response.StatusCode, jsonData)
        return
    }
    //peticion delete al servicio de archivos
    if method == "DELETE" {
        req, err := http.NewRequest("DELETE", "http://10.0.2.4:8082/"+username+"/"+doc_id, nil)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        client := &http.Client{}
        response, err = client.Do(req)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        defer response.Body.Close()
        if response.StatusCode != 204 {
            body, _ := io.ReadAll(response.Body)
            c.JSON(response.StatusCode, gin.H{"error": string(body)})
            return
        }
        c.JSON(response.StatusCode, gin.H{})
        return
    }
    //peticion get de los archivos de un usuario al servicio de archivos
    if method == "LIST" {
        response, err = http.Get("http://10.0.2.4:8082/"+username)
        if err != nil {
            c.JSON(500, gin.H{"error": "could not connect to file service"})
            return
        }
        defer response.Body.Close()
        var jsonFiles FileResponse
        err = json.NewDecoder(response.Body).Decode(&jsonFiles)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(response.StatusCode, gin.H{"files": jsonFiles.Files})
        return
    }

    //validacion del cuerpo del mensaje para los archivos
    valid, json_data, message := checkBody_file(c)
    if !valid {
        c.JSON(400, gin.H{"error": message})
        return
    }
    jsonData, err := json.Marshal(json_data)
    if err != nil {
        c.JSON(500, gin.H{"error": "internal error, json marshal"})
        return
    }
    //peticion post al servicio de archivos
    if method == "POST" {
        response, err = http.Post("http://10.0.2.4:8082/"+username+"/"+doc_id, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            c.JSON(500, gin.H{"error": "could not connect to file service"})
            return
        }
        defer response.Body.Close()
        var data map[string]int
        err = json.NewDecoder(response.Body).Decode(&data)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(response.StatusCode, data)
        return
    }
    //peticion put al servicio de archivos
    if method == "PUT" {
        req, err := http.NewRequest("PUT", "http://10.0.2.4:8082/"+username+"/"+doc_id, bytes.NewBuffer(jsonData))
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        req.Header.Set("Content-Type", "application/json")
        client := &http.Client{}
        response, err = client.Do(req)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        defer response.Body.Close()
        var data map[string]int
        err = json.NewDecoder(response.Body).Decode(&data)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(response.StatusCode, data)
    }
}