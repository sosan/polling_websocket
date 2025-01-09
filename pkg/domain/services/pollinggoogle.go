package services

import "fmt"

func (a *PollingServiceImpl) getAllContentFromGoogleSheets(actionID *string, userID *string, commandType string) (data *string, err error) {
	response, err := a.httpRepo.GetActionByID(actionID, userID, commandType, 1)
	if err != nil {
		return nil, fmt.Errorf("ERROR | failed to fetch action by ID: %w", err)
	}
	// reponse with len 0 can be posible and its NOT ERROR per se
	if len(response.Data) == 0 {
		return nil, nil
	}
	if len(response.Data) > 1 {
		return nil, fmt.Errorf("response contains more than one data entry")
	}

	return &response.Data[0].Data, nil
}
