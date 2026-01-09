package main

// maskPassword masks a password for display in logs and error messages.
// For passwords <= 4 characters, returns "****"
// For longer passwords, shows first 2 and last 2 characters with **** in between
// This prevents password exposure in logs while allowing identification of which credential was used.
//
// Examples:
//   - "abc" -> "****"
//   - "password123" -> "pa****23"
//   - "MySecretP@ss" -> "My****ss"
func maskPassword(password string) string {
	if len(password) <= 4 {
		return "****"
	}
	// Show first 2 and last 2 characters
	return password[:2] + "****" + password[len(password)-2:]
}

// maskUsername masks a username for display in logs and error messages.
// For usernames <= 4 characters, returns "****"
// For longer usernames, shows first 2 and last 2 characters with **** in between
// Useful for email addresses in authentication contexts.
//
// Examples:
//   - "abc" -> "****"
//   - "user@example.com" -> "us****om"
//   - "admin" -> "ad****in"
func maskUsername(username string) string {
	if len(username) <= 4 {
		return "****"
	}
	// Show first 2 and last 2 characters
	return username[:2] + "****" + username[len(username)-2:]
}
