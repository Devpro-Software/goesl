package goesl

import (
	"fmt"
	"strings"
)

type HLClient struct {
	host     string
	port     int
	password string
}

const defaultConnectionTimeout = 10

func NewHLClient(host string, port int, password string) *HLClient {
	return &HLClient{
		host:     host,
		port:     port,
		password: password,
	}
}

func (c *HLClient) makeClient() (*Client, error) {
	return NewClient(c.host, uint(c.port), c.password, defaultConnectionTimeout)
}

func (c *HLClient) Originate(gateway, destination, callerID string, params map[string]string) (string, error) {
	esl, err := c.makeClient()
	if err != nil {
		return "", fmt.Errorf("failed to make client: %w", err)
	}
	go esl.Handle()
	defer esl.Close()

	params["origination_caller_id_number"] = callerID
	paramList := make([]string, 0, len(params))
	for k, v := range params {
		paramList = append(paramList, fmt.Sprintf("%s=%s", k, v))
	}
	paramsStr := strings.Join(paramList, ",")

	err = esl.Api(fmt.Sprintf("originate {%s}sofia/gateway/%s/%s &park()", paramsStr, gateway, destination))
	if err != nil {
		return "", fmt.Errorf("failed to originate call: %w", err)
	}

	var msg *Message
	for {
		m, err := esl.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("failed to read message: %w", err)
		}

		if len(m.Body) > 0 {
			msg = m
			break
		}
	}

	body := string(msg.Body)

	callID := strings.TrimPrefix(body, "+OK ")
	return callID, nil
}

func (c *HLClient) OriginateConference(gateway, destination, callerID, conferenceName string, params map[string]string, conferenceFlags string) (string, error) {
	esl, err := c.makeClient()
	go esl.Handle()
	defer esl.Close()
	if err != nil {
		return "", fmt.Errorf("failed to make client: %w", err)
	}

	params["origination_caller_id_number"] = callerID
	paramList := make([]string, 0, len(params))
	for k, v := range params {
		paramList = append(paramList, fmt.Sprintf("%s=%s", k, v))
	}
	paramsStr := strings.Join(paramList, ",")

	err = esl.Api(fmt.Sprintf("originate {%s}sofia/gateway/%s/%s &conference(%s+flags{%s})", paramsStr, gateway, destination, conferenceName, conferenceFlags))
	if err != nil {
		return "", fmt.Errorf("failed to originate call: %w", err)
	}

	var msg *Message
	for {
		m, err := esl.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("failed to read message: %w", err)
		}

		if len(m.Body) > 0 {
			msg = m
			break
		}
	}
	body := string(msg.Body)
	callID := strings.TrimPrefix(body, "+OK ")
	return callID, nil
}

func (c *HLClient) Bridge(callID1 string, callID2 string) (string, error) {
	esl, err := c.makeClient()
	go esl.Handle()
	defer esl.Close()
	if err != nil {
		return "", fmt.Errorf("failed to make client: %w", err)
	}

	err = esl.Api(fmt.Sprintf("uuid_bridge %s %s", callID1, callID2))
	if err != nil {
		return "", fmt.Errorf("failed to bridge call: %w", err)
	}

	var msg *Message
	for {
		m, err := esl.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("failed to read message: %w", err)
		}

		if len(m.Body) > 0 {
			msg = m
			break
		}
	}

	callID := strings.TrimPrefix(string(msg.Body), "+OK ")
	return callID, nil
}

func (c *HLClient) Deflect(callID, destination, gateway string) (string, error) {
	esl, err := c.makeClient()
	go esl.Handle()
	defer esl.Close()
	if err != nil {
		return "", fmt.Errorf("failed to make client: %w", err)
	}

	err = esl.Api(fmt.Sprintf("uuid_deflect %s sofia/gateway/%s/%s", callID, gateway, destination))
	if err != nil {
		return "", fmt.Errorf("failed to deflect call: %w", err)
	}

	var msg *Message
	for {
		m, err := esl.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("failed to read message: %w", err)
		}

		if len(m.Body) > 0 {
			msg = m
			break
		}
	}
	body := string(msg.Body)
	callID = strings.TrimPrefix(body, "+OK ")
	return callID, nil
}

func (c *HLClient) Transfer(callID, destination string) (string, error) {
	esl, err := c.makeClient()
	go esl.Handle()
	defer esl.Close()
	if err != nil {
		return "", fmt.Errorf("failed to make client: %w", err)
	}

	err = esl.Api(fmt.Sprintf("uuid_transfer %s %s", callID, destination))
	if err != nil {
		return "", fmt.Errorf("failed to transfer call: %w", err)
	}

	var msg *Message
	for {
		m, err := esl.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("failed to read message: %w", err)
		}

		if len(m.Body) > 0 {
			msg = m
			break
		}
	}
	body := string(msg.Body)
	callID = strings.TrimPrefix(body, "+OK ")
	return callID, nil
}
