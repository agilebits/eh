//
// Sample file that is encrypted using AWS KMS keys
//
name = "Sample AWS"

// Special `eh` element that defines encryption parameters
eh {
	// encrypted will changed to `true` after all values are encrypted
	encrypted = false
	key = ""
	
	// Encryption service parameters (AWS)
	service {
		type = "awskms"
		region = "us-east-1"
		masterKey = "arn:aws:kms:us-east-1:123456789012:alias/KeyAlias"
	}

	// List of protected keys. The values of these keys will be encrypted.
	protect = [
		"password",
		"secret",
		"duoApplicationKey",
		"APIKey",
		"emailToDomainHMACSecret",
		"privateKey",
		"securityToken",
		"hook",
	]
}

// Sample data
smtp {
	username = "AKIAI4JG42A2LILVBNNZ"
	password = "smtp-password-value-1"
	host = "email-smtp.us-east-1.amazonaws.com"
	port = 587
}

s3 {
	name = "avatars"
	region = "us-east-1"
	credentials = "accesskey"
	accessKey = "AKIAJZF5WRZKVRMYUIDQ"
	secret = "s3-secret-here"
	securityToken = "security-token-value-here"
	bucket = "some-bucket.com"
	baseURL = "https://s3.amazonaws.com/some-bucket.com/"
}

slack {
	channel1 {
		hook = "https://slack.com/1"
	}

	channel2 {
		hook = "https://slack.com/2"
	}
}
