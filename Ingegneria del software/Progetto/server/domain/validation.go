package domain

// IsValidEAN13 verifica che una stringa sia un barcode EAN-13 valido.
// Controlla sia la forma sia il checksum, così il backend rifiuta codici
// solo apparentemente validi.
func IsValidEAN13(barcode string) bool {
	if len(barcode) != 13 {
		return false
	}

	sum := 0
	for i := 0; i < 12; i++ {
		if barcode[i] < '0' || barcode[i] > '9' {
			return false
		}
		digit := int(barcode[i] - '0')

		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}

	if barcode[12] < '0' || barcode[12] > '9' {
		return false
	}
	checkDigit := int(barcode[12] - '0')

	computed := (10 - (sum % 10)) % 10
	return computed == checkDigit
}
