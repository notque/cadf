package cadf

import (
	"github.com/satori/go.uuid"
	"math/rand"
	"time"
	"strings"
)

//NewEventDetail Generate Fake CADF events for profit
func NewEventDetail() *EventDetail{
	rand.Seed(time.Now().Unix())
	u := uuid.Must(uuid.NewV4())
	uuid := u.String()
	e := EventDetail{
		TypeURI: TypeURI(),
		ID: strings.ToLower(uuid),
		EventTime: time.Now().String(),
		Action: Action(), 
		EventType: EventType(), 
		Outcome: Outcome(), 
	}
	return &e
}

//TypeURI generator
func TypeURI() string {
	s := []string{
		"compute/server",
		"compute/server/volume-attachment", 
		"compute/keypair", 
		"network/floatingip",}
	return RandSlice(s)
}

//Action generator
func Action() string {
	s := []string{
		"create/role_assignment",}
	return RandSlice(s)
}

//EventType generator
func EventType() string {
	s := []string{
		"activity",}
	return RandSlice(s)
}

//Outcome generator
func Outcome() string {
	s := []string{
		"success",
		"failure",}
	return RandSlice(s)
}

//RandSlice Choose random slice entry
func RandSlice(s []string) string {
	n := rand.Int() % len(s)
	return s[n]
}

// EventDetail contains the CADF event according to CADF spec, section 6.6.1 Event (data)
// Extensions: requestPath (OpenStack, IBM), initiator.project_id/domain_id
// Omissions: everything that we do not use or not expose to API users
//  The JSON annotations are for parsing the result from ElasticSearch AND for generating the Hermes API response
type EventDetail struct {
	TypeURI   string `json:"typeURI"`
	ID        string `json:"id"`
	EventTime string `json:"eventTime"`
	Action    string `json:"action"`
	EventType string `json:"eventType"`
	Outcome   string `json:"outcome"`
	Reason    struct {
		ReasonType string `json:"reasonType"`
		ReasonCode string `json:"reasonCode"`
	} `json:"reason,omitempty"`
	Initiator   Resource     `json:"initiator"`
	Target      Resource     `json:"target"`
	Observer    Resource     `json:"observer"`
	Attachments []Attachment `json:"attachments,omitempty"`
	// requestPath is an extension of OpenStack's pycadf which is supported by IBM as well
	RequestPath string `json:"requestPath,omitempty"`
}

// Resource contains attributes describing a (OpenStack-) Resource
type Resource struct {
	TypeURI   string `json:"typeURI"`
	Name      string `json:"name,omitempty"`
	Domain    string `json:"domain,omitempty"`
	ID        string `json:"id"`
	Addresses []struct {
		URL  string `json:"url"`
		Name string `json:"name,omitempty"`
	} `json:"addresses,omitempty"`
	Host *struct {
		ID       string `json:"id,omitempty"`
		Address  string `json:"address,omitempty"`
		Agent    string `json:"agent,omitempty"`
		Platform string `json:"platform,omitempty"`
	} `json:"host,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	// project_id and domain_id are OpenStack extensions (introduced by Keystone and keystone(audit)middleware)
	ProjectID string `json:"project_id,omitempty"`
	DomainID  string `json:"domain_id,omitempty"`
}

// Attachment contains self-describing extensions to the event
type Attachment struct {
	// Note: name is optional in CADF spec. to permit unnamed attachments
	Name string `json:"name,omitempty"`
	// this is messed-up in the spec.: the schema and examples says contentType. But the text often refers to typeURI.
	// Using typeURI would surely be more consistent. OpenStack uses typeURI, IBM supports both
	// (but forgot the name property)
	TypeURI string `json:"typeURI"`
	// Content contains the payload of the attachment. In theory this means any type.
	// In practise we have to decide because otherwise ES does based one first value
	Content string `json:"content"`
}


