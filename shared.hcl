ehcl {
	encrypted = false
	key       = ""

	service {
		type = "local"
	}

	protect = [
		"password",
	]
}

shared_service {
	username = "shared-fragment-user"
	password = "shared-fragment-password"
}