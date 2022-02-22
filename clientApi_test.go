package  clientApi

import (
    "encoding/json"
    "fmt"
    "testing"
)


var urlForAccAPI = "http://192.168.1.9:8080"

func happyCreateFetchDelete(t *testing.T) {

    udd := User_defined_value {
               Key   : "key",
               Value : "value",
    }

    accountInstance :=  AttributesAccount {
                              Country           : "GB",
                              Base_currency     : "GBP",
                              Bank_id           : "400300",
                              Bank_id_code      : "GBDSC",
                              Bic               : "NWBKGB22",
                              Name              : []string { "bob" }, // name is required
                              User_defined_data : []User_defined_value{ udd },
                              Validation_type   :"card",
                              Reference_mask       :"############",
                              Acceptance_qualifier :"same_day",
    }

    localAccount := Account {
                        Type            : "accounts",
                        Id              : "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
                        Organisation_id : "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
                        Attributes      : accountInstance,
    }

    responseStr,status := CreateAccount(urlForAccAPI,localAccount)
    fmt.Print("happyCreateFetchDelete createAccount responseStr ",responseStr)

    // From the response we can extract all data, but for now we care 
    // mostly about the version since that's needed for deletion later.
    var top TopS
    var jsonStr = []byte(responseStr)
    err := json.Unmarshal(jsonStr, &top)
    if err != nil {
			t.Errorf("%v",err)
    } else {
      fmt.Print(" unmarshalled data: ",top)
      fmt.Println(" unmarshalled top.data created_on: ",top.Data.Created_on)
      fmt.Println(" unmarshalled top.data id: ",top.Data.Id)
      fmt.Println(" unmarshalled top.data version: ",top.Data.Version)
    }

		if status != "201 Created" {
			t.Errorf("Positive Create Failed %v",status)
		} else {
			t.Logf("Positive Create Success !")
		}

/*
    request := " {\"data\": " + jsonStr + "}"
    url := urlForAccAPI+"/v1/organisation/accounts"
    fmt.Println("URL:>", url)

    responseStr, status := PostRequest(url,request)

    if status == "201 Created" {
      t.Logf("Creation Succeeded.")
    } else {
      t.Errorf("Create Failed %v",status)
    }
*/


    url := "http://192.168.1.9:8080/v1/organisation/accounts/"
    responseStr, respStatus := GetRecord(url, localAccount.Id)
    // TODO: must parse output for version and store this, needed for delete

    fmt.Println("get record reeponse status:", respStatus)
    _ = responseStr
    if respStatus == "200 OK" {
      t.Logf("Positive Get Succeeded.")
    } else {
      t.Errorf("Positive Get Failed %v",respStatus)
    }


    // delete the record
    respStatus  = DelRecord(url, localAccount.Id, top.Data.Version)
    fmt.Println("delete reeponse status:", respStatus)
    if respStatus == "passed" {
      t.Logf("Positive Delete Succeeded.")
    } else {
      t.Errorf("Positive Delete Failed %v",respStatus)
    }

    // verify record has actually been deleted

    url = "http://192.168.1.9:8080/v1/organisation/accounts/"
    responseStr, respStatus = GetRecord(url, localAccount.Id)
    fmt.Println("get record reeponse status:", respStatus, responseStr)
    if responseStr == "{\"data\":null}" {
      t.Logf("Delete Verification Succeeded.")
    } else {
      // if it exists, then we have failed to delete
      t.Errorf("delete verification failed record %v",respStatus)
    }


}

func sadCreate(t *testing.T) {
// create an invalid record, verify this fails
    udd := User_defined_value {
               Key   : "key",
               Value : "value",
    }

    accountInstance :=  AttributesAccount {
                              Country           : "GB",
                              Base_currency     : "GBP",
                              Bank_id           : "400300",
                              Bank_id_code      : "GBDSC",
                              Bic               : "NWBKGB22",
                              //Name              : []string { "bob" }, // name is required
                              User_defined_data : []User_defined_value{ udd },
                              Validation_type   :"card",
                              Reference_mask       :"############",
                              Acceptance_qualifier :"same_day",
    }

    localAccount := Account {
                        Type            : "accounts",
                        Id              : "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
                        Organisation_id : "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
                        Attributes      : accountInstance,
    }

    responseStr,status := CreateAccount(urlForAccAPI,localAccount)
    fmt.Print("sadCreate responseStr ",responseStr)
		if status == "201 Created" {
			t.Errorf("Negative Create Failed %v",status)
		} else {
			t.Logf("Negative Create Succeeded ")
		}

}


func sadFetch(t *testing.T) {
// fetch a record that does not exist, verify this fails
    id := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
    url := urlForAccAPI + "/v1/organisation/accounts/"
    responseStr, respStatus := GetRecord(url, id)
    fmt.Println("get record reeponse status:", respStatus)
    _ = responseStr
    if respStatus == "200 OK" {
      t.Errorf("Negative Get Failed %v",respStatus)
    } else {
      t.Logf("Negative Get Succeeded.")
    }

}

func sadFetch2(t *testing.T) {
// fetch a record that does not exist, verify this fails
    id := "garbage-data"
    url := urlForAccAPI + "/v1/organisation/accounts/"
    responseStr, respStatus := GetRecord(url, id)
    //fmt.Println("Negative get record reeponse status:", respStatus)
    _ = responseStr
    if respStatus == "200 OK" {
      t.Errorf("Negative Get Failed %v",respStatus)
    } else {
      t.Logf("Negative Get Succeeded.")
    }

}

func happyDelete(t *testing.T) {
    // Create a valid record, verify it exists, then delete it
    // verify the deletion worked.
    udd := User_defined_value {
               Key   : "key",
               Value : "value",
    }

    accountInstance :=  AttributesAccount {
                              Country           : "GB",
                              Base_currency     : "GBP",
                              Bank_id           : "400300",
                              Bank_id_code      : "GBDSC",
                              Bic               : "NWBKGB22",
                              Name              : []string { "bob" }, // name is required
                              User_defined_data : []User_defined_value{ udd },
                              Validation_type   :"card",
                              Reference_mask       :"############",
                              Acceptance_qualifier :"same_day",
    }

    localAccount := Account {
                        Type            : "accounts",
                        Id              : "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
                        Organisation_id : "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
                        Attributes      : accountInstance,
    }

    responseStr,status := CreateAccount(urlForAccAPI,localAccount)
		if status != "201 Created" {
			t.Errorf("setup Failed during create %v",status)
      return
		}
    // From the response we can extract all data, but for now we care 
    // mostly about the version since that's needed for deletion later.
    var top TopS
    var jsonStr = []byte(responseStr)
    err := json.Unmarshal(jsonStr, &top)
    if err != nil {
			t.Errorf("%v",err)
    } else {
      fmt.Println(" unmarshalled top.data version: ",top.Data.Version)
    }



    url := urlForAccAPI + "/v1/organisation/accounts/"
    respStatus := DelRecord(url, localAccount.Id, top.Data.Version)


    fmt.Println("delete response status:", respStatus)
    if respStatus == "passed" {
      t.Logf("Positive Delete Succeeded.")
    } else {
      t.Errorf("Positive Delete Failed.")
    }

}



func TestEverything(t *testing.T) {

  happyCreateFetchDelete(t)
  sadCreate(t)
  sadFetch(t)
  sadFetch2(t)
  happyDelete(t)

}

