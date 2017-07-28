# eh — Encrypted HCL

A small utility to encrypt and decrypt some of the values in .hcl files.

[HashCorp configuration language](https://github.com/hashicorp/hcl) is a great format for application config files. 

Config files often include passwords, private keys, and other secrets. This utility can encrypt these values to protect them. 

## Install

```
go install github.com/agilebits/eh
eh help
```

## Encryption Options

There are two encryption options: "local" and "awskms". 

The local option uses a master key that is stored in `~/.sm/masterkey` file. A new masterkey is created on the first run. It could be shared within the team if the configuration file is checked into the version control system.

For apps running on AWS, the "awskms" option can be used. It is based on the KMS key that should be made available to the EC2 instances.

## Reading Config in Apps

```
    configURL := "file://./config.hcl"
    config, err := secrets.Read(configURL)
	if err != nil {
        ...
	}

    hclObject, err := hcl.ParseBytes(contents)
	if err != nil {
        ...
	}

	var config MyAppConfig
	if err := hcl.DecodeObject(&config, hclObject); err != nil {
		...
	}

```

## Notes

For more complex secret management options, check out [Vault by HashiCorp](https://www.vaultproject.io/) and [Docker Secrets](https://docs.docker.com/engine/swarm/secrets/).