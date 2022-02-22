// A client API to create, fetch, and delete account records on the Form3
// Accounting API
//
// Author: Aly Nathoo
// Date  : Feb 20 2022
//
package clientApi

// TODO: add logging where necessary

import (
    "encoding/json"
//    "fmt"
//    "log"
    "net/http"
    "bytes"
    "io/ioutil"
    "strconv"
)


//////////////////////////////////////////////////////////////////////////////
//   Data Structures
//////////////////////////////////////////////////////////////////////////////

type TopS struct {
  Data     DataStruct
  Links    LinksStruct
}

type User_defined_value struct {
  Key   string
  Value string
}

type LinksStruct struct {
  Self string
}

type DataStruct struct {
  Attributes      AttributesAccount
  Created_on      string
  Modified_on     string
  Type            string
  Id              string
  Organisation_id string
  Version         int
}

type Account struct {
  Type            string
  Id              string
  Organisation_id string
  Attributes      AttributesAccount `json:"attributes,omitempty"`
}

type AttributesAccount struct {
  Alternative_names      []string
  Country                string
  Base_currency          string
  Bank_id                string
  Bank_id_code           string
  Account_number         string
  Bic                    string
  Iban                   string
  Customer_id            string
  Name                   []string
  Account_classification string `json:"account_classification,omitempty"`
  Name_matching_status   string
  User_defined_data      []User_defined_value
  Validation_type        string
  Reference_mask         string
  Acceptance_qualifier   string
  // TODO: Deprecated fields: possible handling: 
  //   1. allow them but give warning if used
  //   2. don't allow any unknown fields, fail immediately
}


//////////////////////////////////////////////////////////////////////////////
//   FUNCTIONS
//////////////////////////////////////////////////////////////////////////////

// pass in account structure
// create an account in local_account structure
// return a json string sequence that can be unmarshalled
func CreateAccount(urlForAccAPI string, local_account  Account) (string, string) {

    var final_string string
    //fmt.Print("\n inside ceateAccount local_account:",local_account)

    b, err := json.Marshal(local_account)
    if err != nil {
      return "","json Marshal failed"
    }
    final_string = string(b)
    //fmt.Println("createAccount")
    //fmt.Println(final_string)
    request := " {\"data\": " + final_string + "}"
    url := urlForAccAPI+"/v1/organisation/accounts"
    //fmt.Println("create URL:", url)

    responseStr, status := PostRequest(url,request)

    return responseStr, status
}


// delete a record, verify it is deleted
// return either "pass" or "fail"
// "fail" implies the record still exists
func DelRecord(url string, id string, version int) (string) {

    reqString := id+"?version="+strconv.Itoa(version)
    //fmt.Print(" delRecord:"+ url + reqString )

    req, err := http.NewRequest("DELETE", url + reqString, nil )

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Since Delete does not provide a clear response, do a get and 
    // verify record no longer exists

    responseStr, respStatus := GetRecord(url, id)
    _ = respStatus

    if responseStr == "{\"data\":null}" {
      return "passed"
    } else {
      // if it exists, then we have failed to delete
      return "failed"
    }
}

// fetch a record with Id
// return a json string sequence that can be unmarshalled
func GetRecord(url string, Id string) (string, string) {

    req, err := http.NewRequest("GET", url, bytes.NewBufferString("/" + Id))
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    if string(body) == "{\"data\":null}" {
       resp.Status = "Failed"
    }

    return string(body),resp.Status

}


func PostRequest(url string, request string) (string, string) {

    req, err := http.NewRequest("POST", url, bytes.NewBufferString(request))

    req.Header.Set("POST", "/v1/organisation/accounts HTTP/1.1")
    req.Header.Add("Content-Type", "application/json")

    //fmt.Println(request)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    //fmt.Println("response Status:", resp.Status)
    //fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    //fmt.Println("response Body:", string(body))
    return string(body),resp.Status
}

//////////////////////////////////////////////////////////////////////////////
/*
curl get example request:
//  curl http://192.168.1.9:8080/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc
curl delete example request:
// curl -X DELETE -H "Content-Type: application/json" http://192.168.1.9:8080/v1/organisation/accounts/"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=0"

 server on create returns:
{  
   "data" : {
      "attributes" : {
         "alternative_names" : null,
         "bank_id" : "400300",
         "bank_id_code" : "GBDSC",
         "base_currency" : "GBP",
         "bic" : "NWBKGB22",
         "country" : "GB",
         "name" : [
            "bob"
         ]
      },
      "created_on" : "2022-02-21T16:06:04.860Z",
      "id" : "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
      "modified_on" : "2022-02-21T16:06:04.860Z",
      "organisation_id" : "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
      "type" : "accounts",
      "version" : 0
   },
   "links" : {
      "self" : "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
   }
}
*/
/*

curl -X POST -H "POST /v1/organisation/accounts HTTP/1.1" -H "Content-Type: application/json" \
    -d '{
  "data": {
{"Type":"accounts","Id":"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc","Organisation_id":"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c","attributes":{"Country":"GB","Base_currency":"GBP","Bank_id":"400300","Bank_id_code":"GBDSC","Account_number":"","Bic":"NWBKGB22","Iban":"","Customer_id":"","Name":null,"Alternative_names":null,"Account_classification":"","Name_matching_status":"","User_defined_data":[{"Key":"key","Value":"value"}],"Validation_type":"card","Reference_mask":"############","Acceptance_qualifier":"same_day"}} } http://localhost:8080/v1/organisation/accounts

curl -X POST -H "POST /v1/organisation/accounts HTTP/1.1" -H "Content-Type: application/json" \
    -d '{
  "data": {
    "type": "accounts",
    "id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
    "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
    "attributes": {
      "country": "GB",
      "base_currency": "GBP",
      "bank_id": "400300",
      "bank_id_code": "GBDSC",
      "bic": "NWBKGB22",
      "name": ["bob"],      
      "user_defined_data": [
        {
          "key": "Some account related key",
          "value": "Some account related value"
        }
      ],
      "validation_type": "card",
      "reference_mask": "############",
      "acceptance_qualifier": "same_day"
    }
  }
}' \
    http://localhost:8080/v1/organisation/accounts


curl -X DELETE -H "Content-Type: application/json" http://localhost:8080/v1/organisation/accounts/"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=0"

*/


