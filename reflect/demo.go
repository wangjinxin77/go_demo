package main
import (
	"fmt"
	"reflect"
)

func main() {
	any_list := []interface{}{1, 3, "wjx"}
	for i, ele := range any_list {
		fmt.Println(reflect.ValueOf(i),":", reflect.TypeOf(ele),"-", reflect.ValueOf(ele))
		fmt.Printf("value %v, %v, %v\n", ele, reflect.TypeOf(ele), reflect.ValueOf(ele))
		fmt.Printf("type %T, %T, %T\n", ele, reflect.TypeOf(ele), reflect.ValueOf(ele))
	}
	str, ok := any_list[2].(string)  // interface{}转string
	fmt.Println("str: ", str, "str_is:", ok)

	// reflect.DeepEqual()
	// 不能使用 == 操作符比较 map，需要使用  reflect.DeepEqual
	// 根据 reflect.DeepEqual 的比较规则，map 的键使用 == 操作符比较， 而值会使用 reflect.DeepEqual  递归比较
	str1, str2 := "wjx", "wjx"
	a := map[string]*string{str1: &str1}
	b := map[string]*string{str2: &str2}
	fmt.Println("reflect.DeepEqual: ", reflect.DeepEqual(a, b))
}
