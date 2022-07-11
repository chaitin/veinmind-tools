package rule

import (
	"github.com/BurntSushi/toml"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/embed"
	"github.com/gobwas/glob"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"regexp"
	"testing"
)

func TestSensitiveRuleUnmarshall(t *testing.T) {
	rules := `
	# ref https://gitlab.com/gitlab-org/security-products/analyzers/secrets/-/blob/master/gitleaks.toml

[whitelist]
paths = ["/usr/**", "/lib/**", "/lib32/**", "/bin/**", "/sbin/**" ,"/var/lib/**", "/var/log/**"]

[[rules]]
id = 1
name = "gitlab_personal_access_token"
description = "GitLab Personal Access Token"
match = '''glpat-[0-9a-zA-Z\-]{20}'''
level = 'high'

[[rules]]
id = 2
name = "AWS"
description = "AWS Access Token"
match = '''AKIA[0-9A-Z]{16}'''
level = 'high'

# Cryptographic keys
[[rules]]
id = 3
name = "PKCS8 private key"
description = "PKCS8 private key"
match = '''-----BEGIN PRIVATE KEY-----'''
level = 'high'

[[rules]]
id = 4
name = "RSA private key"
description = "RSA private key"
match = '''-----BEGIN RSA PRIVATE KEY-----'''
level = 'high'

[[rules]]
id= 5
name = "SSH private key"
description = "SSH private key"
match = '''-----BEGIN OPENSSH PRIVATE KEY-----'''
level = 'high'

[[rules]]
id = 6
name = "PGP private key"
description = "PGP private key"
match = '''-----BEGIN PGP PRIVATE KEY BLOCK-----'''
level = 'high'

[[rules]]
id = 7
name = "Github Personal Access Token"
description = "Github Personal Access Token"
match = '''ghp_[0-9a-zA-Z]{36}'''
level = 'high'

[[rules]]
id = 8
name = "Github OAuth Access Token"
description = "Github OAuth Access Token"
match = '''gho_[0-9a-zA-Z]{36}'''
level = 'high'

[[rules]]
id = 9
name = "SSH (DSA) private key"
description = "SSH (DSA) private key"
match = '''-----BEGIN DSA PRIVATE KEY-----'''
level = 'high'

[[rules]]
id = 10
name = "SSH (EC) private key"
description = "SSH (EC) private key"
match = '''-----BEGIN EC PRIVATE KEY-----'''
level = 'high'

[[rules]]
id = 11
name = "Github App Token"
description = "Github App Token"
match = '''(ghu|ghs)_[0-9a-zA-Z]{36}'''
level = 'high'

[[rules]]
id = 12
name = "Github Refresh Token"
description = "Github Refresh Token"
match = '''ghr_[0-9a-zA-Z]{76}'''
level = 'high'

[[rules]]
id = 13
name = "Shopify shared secret"
description = "Shopify shared secret"
match = '''shpss_[a-fA-F0-9]{32}'''
level = 'high'

[[rules]]
id = 14
name = "Shopify access token"
description = "Shopify access token"
match = '''shpat_[a-fA-F0-9]{32}'''
level = 'high'

[[rules]]
id = 15
name = "Shopify custom app access token"
description = "Shopify custom app access token"
match = '''shpca_[a-fA-F0-9]{32}'''
level = 'high'

[[rules]]
id = 16
name = "Shopify private app access token"
description = "Shopify private app access token"
match = '''shppa_[a-fA-F0-9]{32}'''
level = 'high'

[[rules]]
id = 17
name = "Slack token"
description = "Slack token"
match = '''xox[baprs]-([0-9a-zA-Z]{10,48})?'''
level = 'high'

[[rules]]
id = 18
name = "Stripe"
description = "Stripe"
match = '''(?i)(sk|pk)_(test|live)_[0-9a-z]{10,32}'''
level = 'high'

[[rules]]
id = 19
name = "PyPI upload token"
description = "PyPI upload token"
match = '''pypi-AgEIcHlwaS5vcmc[A-Za-z0-9-_]{50,1000}'''
level = 'high'

[[rules]]
id = 20
name = "Google (GCP) Service-account"
description = "Google (GCP) Service-account"
match = '''\"type\": \"service_account\"'''

[[rules]]
# demo of this match not matching passwords in urls that contain env vars:
# https://match101.com/r/rT9Lv9/3
id = 21
name = "Password in URL"
description = "Password in URL"
match = '''[a-zA-Z]{3,10}:\/\/[^$][^:@\/]{3,20}:[^$][^:@\n\/]{3,40}@.{1,100}'''
level = 'high'

[[rules]]
id = 22
name = "Heroku API Key"
description = "Heroku API Key"
match = '''(?i)(heroku[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12})['\"]'''
level = 'high'

[[rules]]
id = 23
name = "Slack Webhook"
description = "Slack Webhook"
match = '''https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}'''
level = 'medium'

[[rules]]
id = 24
name = "Twilio API Key"
description = "Twilio API Key"
match = '''SK[0-9a-fA-F]{32}'''
level = 'high'

[[rules]]
id = 25
name = "Age secret key"
description = "Age secret key"
match = '''AGE-SECRET-KEY-1[QPZRY9X8GF2TVDW0S3JN54KHCE6MUA7L]{58}'''
level = 'high'

[[rules]]
id = 26
name = "Facebook token"
description = "Facebook token"
match = '''(?i)(facebook[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-f0-9]{32})['\"]'''
level = 'high'

[[rules]]
id = 27
name = "Twitter token"
description = "Twitter token"
match = '''(?i)(twitter[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-f0-9]{35,44})['\"]'''
level = 'high'

[[rules]]
id = 28
name = "Adobe Client ID (Oauth Web)"
description = "Adobe Client ID (Oauth Web)"
match = '''(?i)(adobe[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-f0-9]{32})['\"]'''
level = 'medium'

[[rules]]
id = 29
name = "Adobe Client Secret"
description = "Adobe Client Secret"
match = '''(?i)(p8e-)[a-z0-9]{32}'''
level = 'high'

[[rules]]
id = 30
name = "Alibaba AccessKey ID"
description = "Alibaba AccessKey ID"
match = '''(?i)(LTAI)[a-z0-9]{20}'''
level = 'medium'

[[rules]]
id = 31
name = "Alibaba Secret Key"
description = "Alibaba Secret Key"
match = '''(?i)(alibaba[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9]{30})['\"]'''
level = 'high'

[[rules]]
id = 32
name = "Asana Client ID"
description = "Asana Client ID"
match = '''(?i)(asana[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([0-9]{16})['\"]'''
level = 'medium'

[[rules]]
id = 33
name = "Asana Client Secret"
description = "Asana Client Secret"
match = '''(?i)(asana[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9]{32})['\"]'''
level = 'high'

[[rules]]
id = 34
name = "Atlassian API token"
description = "Atlassian API token"
match = '''(?i)(atlassian[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9]{24})['\"]'''
level = 'high'

[[rules]]
id = 35
name = "Bitbucket client ID"
description = "Bitbucket client ID"
match = '''(?i)(bitbucket[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9]{32})['\"]'''
level = 'medium'

[[rules]]
id = 36
name = "Bitbucket client secret"
description = "Bitbucket client secret"
match = '''(?i)(bitbucket[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9_\-]{64})['\"]'''
level = 'high'

[[rules]]
id = 37
name = "Beamer API token"
description = "Beamer API token"
match = '''(?i)(beamer[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"](b_[a-z0-9=_\-]{44})['\"]'''
level = 'high'

[[rules]]
id = 38
name = "Clojars API token"
description = "Clojars API token"
match = '''(?i)(CLOJARS_)[a-z0-9]{60}'''
level = 'high'

[[rules]]
id = 39
name = "Contentful delivery API token"
description = "Contentful delivery API token"
match = '''(?i)(contentful[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9\-=_]{43})['\"]'''
level = 'high'

[[rules]]
id = 40
name = "Contentful preview API token"
description = "Contentful preview API token"
match = '''(?i)(contentful[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9\-=_]{43})['\"]'''
level = 'high'

[[rules]]
id = 41
name = "Databricks API token"
description = "Databricks API token"
match = '''dapi[a-h0-9]{32}'''
level = 'high'

[[rules]]
id = 42
name = "Discord API key"
description = "Discord API key"
match = '''(?i)(discord[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-h0-9]{64})['\"]'''
level = 'high'

[[rules]]
id = 43
name = "Discord client ID"
description = "Discord client ID"
match = '''(?i)(discord[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([0-9]{18})['\"]'''
level = 'medium'

[[rules]]
id = 44
name = "Discord client secret"
description = "Discord client secret"
match = '''(?i)(discord[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9=_\-]{32})['\"]'''
level = 'high'

[[rules]]
id = 45
name = "Doppler API token"
description = "Doppler API token"
match = '''(?i)['\"](dp\.pt\.)[a-z0-9]{43}['\"]'''
level = 'high'

[[rules]]
id = 46
name = "Dropbox API secret/key"
description = "Dropbox API secret/key"
match = '''(?i)(dropbox[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9]{15})['\"]'''
level = 'high'

[[rules]]
id = 47
name = "Dropbox short lived API token"
description = "Dropbox short lived API token"
match = '''(?i)(dropbox[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"](sl\.[a-z0-9\-=_]{135})['\"]'''
level = 'high'

[[rules]]
id = 48
name = "Dropbox long lived API token"
description = "Dropbox long lived API token"
match = '''(?i)(dropbox[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"][a-z0-9]{11}(AAAAAAAAAA)[a-z0-9\-_=]{43}['\"]'''
level = 'high'

[[rules]]
id = 49
name = "Duffel API token"
description = "Duffel API token"
match = '''(?i)['\"]duffel_(test|live)_[a-z0-9_-]{43}['\"]'''
level = 'high'

[[rules]]
id = 50
name = "Dynatrace API token"
description = "Dynatrace API token"
match = '''(?i)['\"]dt0c01\.[a-z0-9]{24}\.[a-z0-9]{64}['\"]'''
level = 'high'

[[rules]]
id = 51
name = "EasyPost API token"
description = "EasyPost API token"
match = '''(?i)['\"]EZAK[a-z0-9]{54}['\"]'''
level = 'high'

[[rules]]
id = 52
name = "EasyPost test API token"
description = "EasyPost test API token"
match = '''(?i)['\"]EZTK[a-z0-9]{54}['\"]'''
level = 'high'

[[rules]]
id = 53
name = "Fastly API token"
description = "Fastly API token"
match = '''(?i)(fastly[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9\-=_]{32})['\"]'''
level = 'high'

[[rules]]
id = 54
name = "Finicity client secret"
description = "Finicity client secret"
match = '''(?i)(finicity[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9]{20})['\"]'''
level = 'high'

[[rules]]
id = 55
name = "Finicity API token"
description = "Finicity API token"
match = '''(?i)(finicity[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-f0-9]{32})['\"]'''
level = 'high'


[[rules]]
id = 56
name = "Flutterweave public key"
description = "Flutterweave public key"
match = '''(?i)FLWPUBK_TEST-[a-h0-9]{32}-X'''
level = 'medium'

[[rules]]
id = 57
name = "Flutterweave secret key"
description = "Flutterweave secret key"
match = '''(?i)FLWSECK_TEST-[a-h0-9]{32}-X'''
level = 'high'

[[rules]]
id = 58
name = "Flutterweave encrypted key"
description = "Flutterweave encrypted key"
match = '''FLWSECK_TEST[a-h0-9]{12}'''
level = 'high'

[[rules]]
id = 59
name = "Frame.io API token"
description = "Frame.io API token"
match = '''(?i)fio-u-[a-z0-9-_=]{64}'''
level = 'high'

[[rules]]
id = 60
name = "GoCardless API token"
description = "GoCardless API token"
match = '''(?i)['\"]live_[a-z0-9-_=]{40}['\"]'''
level = 'high'

[[rules]]
id = 61
name = "Grafana API token"
description = "Grafana API token"
match = '''(?i)['\"]eyJrIjoi[a-z0-9-_=]{72,92}['\"]'''
level = 'high'

[[rules]]
id = 62
name = "Hashicorp Terraform user/org API token"
description = "Hashicorp Terraform user/org API token"
match = '''(?i)['\"][a-z0-9]{14}\.atlasv1\.[a-z0-9-_=]{60,70}['\"]'''
level = 'high'

[[rules]]
id = 63
name = "Hashicorp Vault batch token"
description = "Hashicorp Vault batch token"
match = '''b\.AAAAAQ[0-9a-zA-Z_-]{156}'''
level = 'high'

[[rules]]
id = 64
name = "Hubspot API token"
description = "Hubspot API token"
match = '''(?i)(hubspot[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-h0-9]{8}-[a-h0-9]{4}-[a-h0-9]{4}-[a-h0-9]{4}-[a-h0-9]{12})['\"]'''
level = 'high'

[[rules]]
id = 65
name = "Intercom API token"
description = "Intercom API token"
match = '''(?i)(intercom[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9=_]{60})['\"]'''
level = 'high'

[[rules]]
id = 66
name = "Intercom client secret/ID"
description = "Intercom client secret/ID"
match = '''(?i)(intercom[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-h0-9]{8}-[a-h0-9]{4}-[a-h0-9]{4}-[a-h0-9]{4}-[a-h0-9]{12})['\"]'''
level = 'high'

[[rules]]
id = 67
name = "Ionic API token"
description = "Ionic API token"
match = '''(?i)ion_[a-z0-9]{42}'''
level = 'high'

[[rules]]
id = 68
name = "Linear API token"
description = "Linear API token"
match = '''(?i)lin_api_[a-z0-9]{40}'''
level = 'high'

[[rules]]
id = 69
name = "Linear client secret/ID"
description = "Linear client secret/ID"
match = '''(?i)(linear[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-f0-9]{32})['\"]'''
level = 'high'

[[rules]]
id = 70
name = "Lob API Key"
description = "Lob API Key"
match = '''(?i)(lob[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]((live|test)_[a-f0-9]{35})['\"]'''
level = 'high'

[[rules]]
id = 71
name = "Lob Publishable API Key"
description = "Lob Publishable API Key"
match = '''(?i)(lob[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]((test|live)_pub_[a-f0-9]{31})['\"]'''
level = 'high'

[[rules]]
id = 72
name = "Mailchimp API key"
description = "Mailchimp API key"
match = '''(?i)(mailchimp[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-f0-9]{32}-us20)['\"]'''
level = 'high'

[[rules]]
id = 73
name = "Mailgun private API token"
description = "Mailgun private API token"
match = '''(?i)(mailgun[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"](key-[a-f0-9]{32})['\"]'''
level = 'high'

[[rules]]
id = 74
name = "Mailgun public validation key"
description = "Mailgun public validation key"
match = '''(?i)(mailgun[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"](pubkey-[a-f0-9]{32})['\"]'''
level = 'high'

[[rules]]
id = 75
name = "Mailgun webhook signing key"
description = "Mailgun webhook signing key"
match = '''(?i)(mailgun[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-h0-9]{32}-[a-h0-9]{8}-[a-h0-9]{8})['\"]'''
level = 'high'

[[rules]]
id = 76
name = "Mapbox API token"
description = "Mapbox API token"
match = '''(?i)(pk\.[a-z0-9]{60}\.[a-z0-9]{22})'''
level = 'high'

[[rules]]
id = 77
name = "messagebird-api-token"
description = "MessageBird API token"
match = '''(?i)(messagebird[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9]{25})['\"]'''
level = 'high'

[[rules]]
id = 78
name = "MessageBird API client ID"
description = "MessageBird API client ID"
match = '''(?i)(messagebird[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-h0-9]{8}-[a-h0-9]{4}-[a-h0-9]{4}-[a-h0-9]{4}-[a-h0-9]{12})['\"]'''
level = 'medium'

[[rules]]
id = 79
name = "New Relic user API Key"
description = "New Relic user API Key"
match = '''['\"](NRAK-[A-Z0-9]{27})['\"]'''
level = 'high'

[[rules]]
id = 80
name = "New Relic user API ID"
description = "New Relic user API ID"
match = '''(?i)(newrelic[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([A-Z0-9]{64})['\"]'''
level = 'medium'

[[rules]]
id = 81
name = "New Relic ingest browser API token"
description = "New Relic ingest browser API token"
match = '''['\"](NRJS-[a-f0-9]{19})['\"]'''
level = 'high'

[[rules]]
id = 82
name = "npm access token"
description = "npm access token"
match = '''(?i)['\"](npm_[a-z0-9]{36})['\"]'''
level = 'high'

[[rules]]
id = 83
name = "Planetscale password"
description = "Planetscale password"
match = '''(?i)pscale_pw_[a-z0-9\-_\.]{43}'''
level = 'high'

[[rules]]
id = 84
name = "Planetscale API token"
description = "Planetscale API token"
match = '''(?i)pscale_tkn_[a-z0-9\-_\.]{43}'''
level = 'high'

[[rules]]
id = 85
name = "Postman API token"
description = "Postman API token"
match = '''(?i)PMAK-[a-f0-9]{24}\-[a-f0-9]{34}'''
level = 'high'

[[rules]]
id = 86
name = "Pulumi API token"
description = "Pulumi API token"
match = '''pul-[a-f0-9]{40}'''
level = 'high'

[[rules]]
id = 87
name = "Rubygem API token"
description = "Rubygem API token"
match = '''rubygems_[a-f0-9]{48}'''
level = 'high'

[[rules]]
id = 88
name = "Sendgrid API token"
description = "Sendgrid API token"
match = '''(?i)SG\.[a-z0-9_\-\.]{66}'''
level = 'high'

[[rules]]
id = 89
name = "Sendinblue API token"
description = "Sendinblue API token"
match = '''(?i)xkeysib-[a-f0-9]{64}\-[a-z0-9]{16}'''
level = 'high'

[[rules]]
id = 90
name = "Shippo API token"
description = "Shippo API token"
match = '''shippo_(live|test)_[a-f0-9]{40}'''
level = 'high'

[[rules]]
id = 91
name = "Linkedin Client secret"
description = "Linkedin Client secret"
match = '''(?i)(linkedin[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z]{16})['\"]'''
level = 'high'

[[rules]]
id = 92
name = "Linkedin Client ID"
description = "Linkedin Client ID"
match = '''(?i)(linkedin[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9]{14})['\"]'''
level = 'medium'

[[rules]]
id = 93
name = "Twitch API token"
description = "Twitch API token"
match = '''(?i)(twitch[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}['\"]([a-z0-9]{30})['\"]'''
level = 'high'

[[rules]]
id = 94
name = "Typeform API token"
description = "Typeform API token"
match = '''(?i)(typeform[a-z0-9_ .\-,]{0,25})(=|>|:=|\|\|:|<=|=>|:).{0,5}(tfp_[a-z0-9\-_\.=]{59})'''
level = 'high'

[[rules]]
id = 95
name = "Social Security Number"
description = "Social Security Number"
match = '''\d{3}-\d{2}-\d{4}'''
level = 'low'

[[rules]]
id = 96
name = "Version Control File"
description = "Version Control File"
filepath = '''.*\/\.(git|svn)$'''
level = 'high'

[[rules]]
id = 97
name = "Config File"
description = "Config File"
filepath = '''.*\/config\.ini$'''
level = 'medium'

[[rules]]
id = 98
name = "Password In Enviroment"
description = "Password In Enviroment"
env = '''[^=]*(auth|pass|key|secret|token)[^=]*=.*'''
level = 'high'

[[rules]]
id = 99
name = "Desktop Services Store"
description = "Desktop Services Store"
filepath = ''' .*\/\.DS_Store$'''
level = 'low'

[[rules]]
id = 100
name = "MySQL client command history file"
description = "MySQL client command history file"
filepath = '''.*\/\.(mysql|psql|irb)_history$'''
level = 'low'

[[rules]]
id = 101
name = "Recon-ng web reconnaissance framework API key database"
description = "Recon-ng web reconnaissance framework API key database"
filepath = '''.*\/\.recon-ng\/keys\.db$'''
level = 'medium'

[[rules]]
id = 102
name = "DBeaver SQL database manager configuration file"
description = "DBeaver SQL database manager configuration file"
filepath = '''.*\/\.dbeaver-data-sources\.xml$'''
level = 'low'

[[rules]]
id = 103
name = "S3cmd configuration file"
description = "S3cmd configuration file"
filepath = '''.*\/\.s3cfg$'''
level = 'low'

[[rules]]
id = 104
name = "Ruby On Rails secret token configuration file"
description = "If the Rails secret token is known, it can allow for remote code execution. (http://www.exploit-db.com/exploits/27527/)"
filepath = '''.*\/secret_token\.rb$'''
level = 'high'

[[rules]]
id = 105
name = "OmniAuth configuration file"
description = "The OmniAuth configuration file might contain client application secrets."
filepath = '''.*\/omniauth\.rb$'''
level = 'high'

[[rules]]
id = 106
name = "Carrierwave configuration file"
description = "Can contain credentials for online storage systems such as Amazon S3 and Google Storage."
filepath = '''.*\/carrierwave\.rb$'''
level = 'high'

[[rules]]
id = 107
name = "Potential Ruby On Rails database configuration file"
description = "Might contain database credentials."
filepath = '''.*\/database\.yml$'''
level = 'high'

[[rules]]
id = 108
name = "Django configuration file"
description = "Might contain database credentials, online storage system credentials, secret keys, etc."
filepath = '''.*\/settings\.py$'''
level = 'low'

[[rules]]
id = 109
name = "PHP configuration file"
description = "Might contain credentials and keys."
filepath = '''.*\/config(\.inc)?\.php$'''
level = 'low'

[[rules]]
id = 110
name = "Jenkins publish over SSH plugin file"
description = "Jenkins publish over SSH plugin file"
filepath = '''.*\/jenkins\.plugins\.publish_over_ssh\.BapSshPublisherPlugin\.xml$'''
level = 'high'

[[rules]]
id = 111
name = "Potential Jenkins credentials file"
description = "Potential Jenkins credentials file"
filepath = '''.*\/credentials\.xml$'''
level = 'high'

[[rules]]
id = 112
name = "Apache htpasswd file"
description = "Apache htpasswd file"
filepath = '''.*\/\.htpasswd$'''
level = 'low'

[[rules]]
id = 113
name = "Configuration file for auto-login process"
description = "Might contain username and password."
filepath = '''.*\/\.(netrc|git-credentials)$'''
level = 'high'

[[rules]]
id = 114
name = "Potential MediaWiki configuration file"
description = "Potential MediaWiki configuration file"
filepath = '''.*\/LocalSettings\.php$'''
level = 'high'

[[rules]]
id = 115
name = "Rubygems credentials file"
description = "Might contain API key for a rubygems.org account."
filepath = '''.*\/\.gem\/credentials$'''
level = 'high'

[[rules]]
id = 116
name = "Potential MSBuild publish profile"
description = "Potential MSBuild publish profile"
filepath = '''.*\/\.pubxml(\.user)?$'''
level = 'low'
	`
	rulesE := SensitiveConfig{}
	err := toml.Unmarshal([]byte(rules), &rulesE)
	if err != nil {
		t.Error(err)
	}

	expected := SensitiveRule{
		Id:          1,
		Name:        "gitlab_personal_access_token",
		Description: "GitLab Personal Access Token",
		Level:       "high",
		Match:       `glpat-[0-9a-zA-Z\-]{20}`,
	}

	assert.Equal(t, expected, rulesE.Rules[0])
}

func TestSensitiveRuleUnmarshall2(t *testing.T) {
	rules, err := embed.FS.Open("rules.toml")
	if err != nil {
		t.Error(err)
	}

	rulesByte, err := ioutil.ReadAll(rules)
	if err != nil {
		t.Error(err)
	}

	rulesE := SensitiveConfig{}
	err = toml.Unmarshal(rulesByte, &rulesE)
	if err != nil {
		t.Error(err)
	}

	expected := SensitiveRule{
		Id:          1,
		Name:        "gitlab_personal_access_token",
		Description: "GitLab Personal Access Token",
		Level:       "high",
		Match:       `glpat-[0-9a-zA-Z\-]{20}`,
	}

	expected2 := SensitiveWhiteList{
		Paths: []string{"/usr/**", "/lib/**", "/lib32/**", "/bin/**", "/sbin/**", "/var/lib/**", "/var/log/**"},
	}

	assert.Equal(t, expected, rulesE.Rules[0])
	assert.Equal(t, expected2, rulesE.WhiteList)
}

func TestRule2(t *testing.T) {
	m, err := regexp.MatchString(`.*\/\.(git|svn)$`, "/var/www/html/.git")
	if err != nil {
		t.Error(err)
	}

	assert.True(t, m)
}

func TestRule3(t *testing.T) {
	g, err := glob.Compile(`/etc/ImageMagick-6/mime.xml`)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, g.Match("/etc/ImageMagick-6/mime.xml"), true)
}
