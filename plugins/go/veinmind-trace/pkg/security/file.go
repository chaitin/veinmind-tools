package security

import (
	"io/fs"
	"regexp"
)

var SensitiveDirPerm = map[string]Perm{
	"/etc/shadow":        {0, 0640},
	"/etc/passwd":        {0, 0644},
	"/etc/group":         {0, 0644},
	"/etc/gshadow":       {0, 0640},
	"/etc/ld.so.conf":    {0, 0644},
	"/etc/hosts":         {0, 0644},
	"/etc/hosts.allow":   {0, 0644},
	"/etc/sudoers":       {0, 0640},
	"/etc/ld.so.preload": {0, 0600},
	"/lib/x86_64-linux-gnu/security/pam_unix.so": {0, 644},

	// crontab
	"/var/spool/cron/crontabs":      {0, fs.ModeDir | 0755},
	"/var/spool/cron/crontabs/root": {0, 0600},
	"/etc/crontab":                  {0, fs.ModeDir | 0644},
	"/etc/cron.d":                   {0, fs.ModeDir | 0755},
	"/etc/cron.daily":               {0, fs.ModeDir | 0755},
	"/etc/cron.hourly":              {0, fs.ModeDir | 0755},
	"/etc/cron.monthly":             {0, fs.ModeDir | 0755},
	"/etc/cron.weekly":              {0, fs.ModeDir | 0755},

	"/bin/":   {0, 0},
	"/sbin/":  {0, 0},
	"/lib/":   {0, 0},
	"/lib64/": {0, 0},
	"/usr/":   {0, 0},
	"/run/":   {0, 0},
	"/proc/":  {0, 0},
	"/root":   {0, fs.ModeDir | 0700},
}

type Perm struct {
	Uid  uint32
	Mode fs.FileMode
}

var CDKTrace = []*regexp.Regexp{
	// exp/mount_device.go
	regexp.MustCompile(`/tmp/cdk_.*?`),
	// exp/mount_cgroup.go

	regexp.MustCompile(`/tmp/cgrp/cdk/notify_on_release`),
	regexp.MustCompile(`cdk_cgres_.*?`),
	// exp/rewrite_cgroup_devices.go
	regexp.MustCompile(`/tmp/cdk_dcgroup.*?`),

	// exp/mount-procfs.go
	regexp.MustCompile(`/mnt/host_proc`),
	// containerd_shim_pwn.go

	// exp/abuse_unpriv_userns.go
	regexp.MustCompile(`/mnt/cgrp1`),

	// exp/lxcfs_rw_cgroup.go
	regexp.MustCompile(`/cdk_cgexp_.*?\.sh`),
}
