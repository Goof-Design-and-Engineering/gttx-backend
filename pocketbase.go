// pocketbase.go
package main

import (
	"fmt"
	"log"
	"net/mail"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

// type DNSCreate struct {
// 	Secretapikey string `json:"secretapikey"`
// 	Apikey       string `json:"apikey"`
// 	Name         string `json:"name"`
// 	Type         string `json:"type"`
// 	Content      string `json:"content"`
// 	TTL          string `json:"ttl"`
// }

// const (
// 	ROOTNAME = "gttx.app"
// 	APPNAME  = "gttx.app"
// )

// func addOrgToDNS(orgName string) {
// 	// post to porkbun
// 	httpposturl := fmt.Sprintf("https://porkbun.com/api/json/v3/dns/create/%s", ROOTNAME)
// 	log.Println("HTTP JSON POST URL:", httpposturl)

// 	// read values
// 	data := DNSCreate{
// 		Secretapikey: os.Getenv("PORKBUN_PRIVATE"),
// 		Apikey:       os.Getenv("PORKBUN_PUBLIC"),
// 		Name:         orgName,
// 		Type:         "ALIAS",
// 		Content:      APPNAME,
// 		TTL:          "600",
// 	}
// 	jsonData, _ := json.Marshal(data)

// 	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
// 	if error != nil {
// 		log.Printf("Error processing DNS records: \n%v", error)
// 	}
// 	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

// 	client := &http.Client{}
// 	response, error := client.Do(request)
// 	if error != nil {
// 		log.Printf("Error processing DNS records: \n%v", error)
// 	}
// 	defer response.Body.Close()

// 	log.Println("response Status:", response.Status)
// 	log.Println("response Headers:", response.Header)
// 	body, _ := ioutil.ReadAll(response.Body)
// 	log.Println("response Body:", string(body))

// }

func main() {
	app := pocketbase.New()

	// app.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
	// 	log.Println(e.Record) // still unsaved
	// 	log.Println(e.Record.Collection().TableName())
	// 	log.Println(e.Record.Collection().Name)
	// 	if e.Record.Collection().Name == "organization" {
	// 		if val, ok := e.Record.SchemaData()["name"]; ok {
	// 			result := fmt.Sprintf("%s", val)
	// 			log.Print(result)
	// 			addOrgToDNS(result)

	// 		}
	// 	}
	// 	return nil
	// })

	// app.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
	// 	if e.Record.Collection().Name == "invites" {
	// if users, ok := e.Record.SchemaData()["users"]; ok {
	// 			if reflect.TypeOf(users) == []string {
	// 				for user := range users.([]string) {

	// 					}
	// 			}

	// 			// var typed_users []any
	// 		}
	// 	}
	// })

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/api/send_invite_code", func(c echo.Context) error {
			// or to get the authenticated record:
			// authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			// if authRecord == nil {
			// 	return apis.NewForbiddenError("Only auth records can access this endpoint", nil)
			// }
			log.Println("ENTERED THING")
			log.Printf("RECORD(1):\t%v\n", c.Get(apis.ContextAuthRecordKey).(*models.Record))

			if collection := c.Get(apis.ContextAuthRecordKey).(*models.Record).SchemaData(); collection != nil {
				fmt.Printf("RECORD(2):\t%v\n", collection)
				// verify that the user exists
				role, ok := collection["role"]
				if !ok {
					return c.String(403, "Not logged in")
				}

				if role != "facilitator" {
					return c.String(403, "Not the right role")
				}

				email, ok := collection["email"]
				if !ok {
					return c.String(403, "No email set")
				}

				log.Println(c.QueryParams())
				// print the data sent over
				log.Printf("RECIEVED EMAIL: %s\n", c.QueryParam("email"))
				log.Printf("RECIEVED ()_code: %s\n", c.QueryParam("participant_code"))
				log.Printf("RECIEVED ()_code: %s\n", c.QueryParam("facilitator_code"))
				log.Printf("RECIEVED ()_code: %s\n", c.QueryParam("observer_code"))
				log.Printf("RECIEVED DESIRED_ROLE: %s\n", c.QueryParam("desired_role"))

				codeToSend := ""
				switch c.QueryParam("desired_role") {
				case "facilitator":
					codeToSend = c.QueryParam("facilitator_code")
				case "observer":
					codeToSend = c.QueryParam("observer_code")
				case "participant":
					codeToSend = c.QueryParam("participant_code")
				default:
					return c.String(400, "No `desired_role` set")
				}

				if codeToSend == "" {
					return c.String(400, fmt.Sprintf("The `desired_role`=%s does not have an invite code set", c.QueryParam("desired_role")))
				}

				message := &mailer.Message{
					From: mail.Address{
						Address: app.Settings().Meta.SenderAddress,
						Name:    app.Settings().Meta.SenderName,
					},
					To:      mail.Address{Address: email.(string)},
					Subject: "Welcome to GTTX",
					HTML:    fmt.Sprintf("Welcome!\nhere is your invite code:%s. Goto %s/signup to login and put in the code sent", codeToSend, app.Settings().Meta.AppUrl),
					// bcc, cc, attachments and custom headers are also supported...
					Cc: []string{"kjo018@latech.edu"},
				}

				err := app.NewMailClient().Send(message)
				if err != nil {
					return c.String(403, fmt.Sprintf("Got error  %v", err))

				} else {
					return c.String(200, "successfully sent email")

				}

			}
			return c.String(403, ("Got unknown error oops"))
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
