package file

import (
	"io/fs"
)

var sensitiveDirPerm = map[string]Perm{
	"/etc/shadow":                   {0, 0640},
	"/etc/passwd":                   {0, 0644},
	"/etc/group":                    {0, 0644},
	"/etc/gshadow":                  {0, 0640},
	"/etc/ld.so.conf":               {0, 0644},
	"/etc/hosts":                    {0, 0644},
	"/etc/hosts.allow":              {0, 0644},
	"/etc/sudoers":                  {0, 0640},
	"/etc/ld.so.preload":            {0, 0600},
	"/var/spool/cron/crontabs":      {0, 0730},
	"/var/spool/cron/crontabs/root": {0, 0600},
	"/lib/x86_64-linux-gnu/security/pam_unix.so": {0, 644},
	"/bin/":   {0, 0},
	"/sbin/":  {0, 0},
	"/lib/":   {0, 0},
	"/lib64/": {0, 0},
	"/usr/":   {0, 0},
	"/run/":   {0, 0},
	"/proc/":  {0, 0},
	"/root":   {0, 0750},
}

type Perm struct {
	uid  uint32
	mode fs.FileMode
}
