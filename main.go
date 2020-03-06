package main

import (
	"fmt"
	"github.com/kradalby/opnsense-go/opnsense"
	//"github.com/satori/go.uuid"
	"log"
)

func main() {
	c, err := opnsense.NewClient("http://192.168.45.128",
		"98JfUoQ9QI8YUsUZ6h1eGG9PHaEFnBELf5ANlV/lTv5bWchyPOyZVV2a70rdDV0PvThBwzYYSLR0c/J1",
		"MlMrlTkepR+CNuqpZ/dpBYcloiWogKPtgixM988MuPK7hYUtMFUsI1AQoHhfuqhEdEDPK0BXW8SN7V/0",
		true)

	if err != nil {
		log.Fatal(err)
	}

	// TEST ALIAS
	/*var default_name = "LUUL"
	var default_desc = default_name

	// get all alias
	var all_uiid, _ = c.AliasGetUUIDs()
	if len(all_uiid) > 0 {
		fmt.Println(all_uiid[0])

		// get info for the first alias in the list
		var u = all_uiid[0]
		var resp_get, _ = c.AliasGet(*u)
		fmt.Println("name :", resp_get.Name)
		fmt.Println("uiid :", resp_get.UUID)
		fmt.Println("description :", resp_get.Description)
		fmt.Println("enabled :", resp_get.Enabled)
		fmt.Println("type :", resp_get.Type)
		fmt.Println("content :", resp_get.Content)
		fmt.Println("proto :", resp_get.Proto)

		default_name = resp_get.Name
		default_desc = resp_get.Description

		// delete this alias
		var resp_del, _ = c.AliasDelete(*u)
		fmt.Println("DEL :", resp_del)
	}

	// create or recreate the same alias
	var add_set opnsense.AliasSet
	add_set.Description = default_desc + " NEW "
	add_set.Name = default_name
	add_set.Enabled = "1"
	add_set.Type = "host"
	add_set.Proto = ""
	add_set.Updatefreq = ""
	add_set.Content = ""
	add_set.Counters = ""
	var resp_add, err_add = c.AliasAdd(add_set)
	if err_add != nil {
		log.Println(err_add)
	}

	fmt.Println("new uiid added : ", resp_add)

	var resp_get_after_add, _ = c.AliasGet(*resp_add)

	// update the created alias
	var update_set opnsense.AliasSet
	update_set.Description = resp_get_after_add.Description + " UPDATE "
	update_set.Name = resp_get_after_add.Name
	update_set.Enabled = "1"
	update_set.Type = "host"
	update_set.Proto = ""
	update_set.Updatefreq = ""
	update_set.Content = "178.132\nLLLLLL"
	update_set.Counters = ""
	var _, err_update = c.AliasUpdate(*resp_add, update_set)

	if err_update != nil {
		log.Println(err_update)
	}

	// TEST ALIAS PUSH
	/*alias_name := "TEST"
	var add_set opnsense.AliasPushSet
	add_set.Address = "10.0.0.11"
	var res_add, err_add = c.AliasPushAdd(alias_name, add_set)

	if err_add != nil {
		log.Println(err_add)
		fmt.Println(res_add)
		return
	}

	var add_del opnsense.AliasPushSet
	add_del.Address = "10.0.0.10"
	var res_del, err_del = c.AliasPushDel(alias_name, add_del)

	if err_add != nil {
		log.Println(err_del)
		fmt.Println(res_del)
	}*/

	/*uuid, err_uuid := uuid.FromString("fd89c801-adb3-42f5-94d7-293655ac205c")
	if err_uuid != nil {
		fmt.Printf("Something went wrong: %s", err_uuid)
		return
	}
	alias, err_get := c.AliasGet(uuid)
	if err_get != nil {
		fmt.Printf("Something went wrong: %s", err_get)
		return
	}

	alias.Content = append(alias.Content, "10.0.0.1")

	_, err_update := c.AliasUpdate(uuid, *alias)

	if err_update != nil {
		log.Println(err_update)
	}*/

	list, err_list := c.AliasGetList()
	if err_list != nil {
		log.Println(err_list)
	}

	fmt.Printf("%v", *list)
}
