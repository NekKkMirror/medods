package util

import "fmt"

// SendSecurityAlertEmail is mock implementation for sending email
func SendSecurityAlertEmail(email, ipAddress string) {
	fmt.Printf("Sent security alert to %s about IP change to %s\n", email, ipAddress)
}
