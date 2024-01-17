# Hetzner Dynamic DNS Configurator
Configure a dns records in Hetzner DNS with your config file

# Features
* Support Dynamic DNS for your Records
* Get your IP Address and set as Record
* Configure Records from config File
* Works with Vault

# Configuration

### Sample Config

```yaml
zones:
  - name: your.domain
    ttl: 3600
    records:
      - name: cloud
        type: A
        value: 1.1.1.1
        ttl: 60
      - name: vpn
        type: A
        ttl: 60
      - name: dev
        type: A
        ttl: 60
```

When `value` is specified, then set this address to record, if empty set as address you public ip address

### Environment variables
| ENV VARIABLE                            | DEFAULT | DESCRIPTION                      |
|-----------------------------------------|---------|----------------------------------|
| HETZNER_DNS_VAULT_ADDR               | not set | Vault addr https://vault.local   |
| HETZNER_DNS_VAULT_AUTH_NAME             | not set | Vault approle auth name          |         
| HETZNER_DNS_VAULT_ROLE_ID         | not set | Vault AppRole ID                 | 
| HETZNER_DNS_VAULT_SECRET_STORE_NAME               | not set | Vault secret store name          |
| HETZNER_DNS_VAULT_SECRET_STORE_PATH      | not set | Vault Secret Path                | 
| HETZNER_DNS_ZONES_CONFIG_FILE       | not set | Path to config file in yaml      |  
| HETZNER_DNS_ZONES_UPDATE_INTERVAL               | not set | Time for records update interval |  
| HETZNER_DNS_USE_VAULT              | not set | Use Vault true/false             |  
| HETZNER_DNS_TOKEN         | not set | Token if no use Vault            |  


### Sample .env

```shell
HETZNER_DNS_VAULT_ADDR=https://vault.local.domain:8080
HETZNER_DNS_VAULT_AUTH_NAME=approle
HETZNER_DNS_VAULT_ROLE_ID=8sa45a11-ec85-2cf5-7d95-b487dfas96
HETZNER_DNS_VAULT_SECRET_STORE_NAME=secrets
HETZNER_DNS_VAULT_SECRET_STORE_PATH=hetzner/dns
HETZNER_DNS_ZONES_CONFIG_FILE=./cfg/your.domain.zones.yaml
HETZNER_DNS_ZONES_UPDATE_INTERVAL=60
HETZNER_DNS_USE_VAULT=true
```
