package main
import (
	"fmt"
	"github.com/google/uuid"
)

func UUIDTest(){
	newUuid := uuid.New()  // new uuid
	uuidStr := newUuid.String()  // uuid => string
	fmt.Println("newUuid ", uuidStr)
	fmt.Printf("type: %T, value: %v \n", uuidStr, uuidStr)
	oldUuid, err := uuid.Parse(uuidStr)  //string => uuid
	if err != nil {
		return
	}
	fmt.Printf("oldUuid type: %T, value: %v\n", oldUuid, oldUuid.String())
}

func main() {
	UUIDTest()
}
