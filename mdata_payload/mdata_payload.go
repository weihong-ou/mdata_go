package mdata_payload

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"strings"
)

type MdPayload struct {
	Action string
	Gtin   string
	Mtrl   string
}

func FromBytes(payloadData []byte) (*MdPayload, error) {
	if payloadData == nil {
		return nil, &processor.InvalidTransactionError{Msg: "Must contain payload"}
	}

	parts := strings.Split(string(payloadData), ",")
	if len(parts) != 3 {
		return nil, &processor.InvalidTransactionError{Msg: "Payload is malformed"}
	}

	payload := MdPayload{}
	payload.Action = parts[0]
	payload.Gtin = parts[1]

	if len(payload.Action) < 1 {
		return nil, &processor.InvalidTransactionError{Msg: "Action is required"}
	}

	if len(payload.Gtin) < 1 {
		return nil, &processor.InvalidTransactionError{Msg: "Gtin is required"}
	}

	if payload.Action == "create" || payload.Action == "update" {
		payload.Mtrl = parts[2]
		if len(payload.Mtrl) < 1 {
			return nil, &processor.InvalidTransactionError{Msg: "Mtrl is required for create and update"}
		}
	}

	if strings.Contains(payload.Gtin, "|") {
		return nil, &processor.InvalidTransactionError{
			Msg: fmt.Sprintf("Invalid Name (char '|' not allowed): '%v'", parts[1])}
	}

	return &payload, nil
}
