// Pacakge domain containing domain libraries
package domain

const (
  // ReqKey is the key constant for RequestID
  ReqKey string = "RequestID"
  // PetID is the key constant for RequestID
  PetID string = "PetID"
  // PetsTable is the table in DynamoDB
  PetsTable = "furriend_pets"
)

// Pet is the domain
type Pet struct {
	PetID string
	Name  string
	Type  string
}