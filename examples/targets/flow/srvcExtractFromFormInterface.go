package main

// ExtractFormInterface extract from form interface
func ExtractFormInterface(formInterface interface{}, key string) string {
	// get text
	form := formInterface.(map[string]interface{})
	text := ""
	if rawTexts, exists := form[key].([]interface{}); exists {
		rawText := rawTexts[0]
		text = rawText.(string)
	}
	return text
}
