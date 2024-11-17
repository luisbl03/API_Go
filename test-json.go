package mai
import (
	"fmt"
	"github.com/luideiz/API_Go/repository"
	"github.com/luideiz/API_Go/models"
)

func main() {
	var user models.User
	models.SetUsername(&user, "luideiz")
	models.SetPassword(&user, "1234")
	status := repository.Add(user)
	fmt.Println(status)
}