package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func ValidateUserExists(userID uuid.UUID) error {
    userServiceURL := os.Getenv("USER_SERVICE_URL")
    resp, err := http.Get(fmt.Sprintf("%s/users?id=%s", userServiceURL, userID))
    if err != nil {
        return fmt.Errorf("failed to validate user: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("user not found or inactive")
    }

    var users []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
        return fmt.Errorf("failed to decode user response: %w", err)
    }

    if len(users) == 0 {
        return fmt.Errorf("user not found")
    }

    return nil
}

func ValidatePlanExists(planID uuid.UUID) error {
    planServiceURL := os.Getenv("PLAN_SERVICE_URL")
    resp, err := http.Get(fmt.Sprintf("%s/plans?id=%s", planServiceURL, planID))
    if err != nil {
        return fmt.Errorf("failed to validate plan: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("plan not found or inactive")
    }

    var plans []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&plans); err != nil {
        return fmt.Errorf("failed to decode plan response: %w", err)
    }

    if len(plans) == 0 {
        return fmt.Errorf("plan not found")
    }

    return nil
}