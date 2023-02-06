package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/group"
	"github.com/chaitin/veinmind-common-go/passwd"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"github.com/chaitin/veinmind-common-go/service/report/service"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-basic/pkg/capability"
	"github.com/opencontainers/runtime-spec/specs-go"
)

var (
	reportService = &report.Service{}
	rootCommand   = &cmd.Command{}
	scanCommand   = &cmd.Command{
		Use:   "scan",
		Short: "scan mode",
	}
	scanImageCommand = &cmd.Command{
		Use:   "image",
		Short: "scan image basic info",
	}
	scanContainerCommand = &cmd.Command{
		Use:   "container",
		Short: "scan container basic info",
	}
)

func scanImage(c *cmd.Command, image api.Image) error {
	refs, err := image.RepoRefs()
	if err != nil {
		// no reference image will report ans use sha256 fill repo field
		log.Error(err)
	}

	oci, err := image.OCISpecV1()
	if err != nil {
		return err
	}

	evt := &event.Event{
		BasicInfo: &event.BasicInfo{
			ID:         image.ID(),
			Object:     event.NewObject(image),
			Time:       time.Now(),
			Level:      event.None,
			DetectType: event.Image,
			AlertType:  event.BasicImage,
			EventType:  event.Info,
		},
		DetailInfo: &event.DetailInfo{
			AlertDetail: &event.ImageBasicDetail{
				References:  refs,
				CreatedTime: oci.Created.Unix(),
				Env:         oci.Config.Env,
				Entrypoint:  oci.Config.Entrypoint,
				Cmd:         oci.Config.Cmd,
				WorkingDir:  oci.Config.WorkingDir,
				Author:      oci.Author,
			},
		},
	}

	err = reportService.Client.Report(evt)
	if err != nil {
		return err
	}

	return nil
}

func scanContainer(c *cmd.Command, container api.Container) error {
	var (
		containerRuntime  event.RuntimeType
		rootProcessDetail event.RootProcessDetail
		mountDetails      []event.MountDetail
		processDetails    []event.ProcessDetail
		createdTime       int64
		runtimeUniqDesc   string
	)

	ocispec, err := container.OCISpec()
	if err != nil {
		// if container not running, doesn't exist oci
		ocispec = &specs.Spec{}
	}

	ocistate, err := container.OCIState()
	if err != nil {
		ocistate = &specs.State{}
	}

	switch c := container.(type) {
	case *docker.Container:
		// runtime type
		containerRuntime = event.Docker

		// runtime desc
		runtimeUniqDesc = c.Runtime().UniqueDesc()

		// docker config
		config, err := c.Config()
		if err != nil {
			log.Error(err)
		} else {
			// created time
			createdTime = config.Created.Unix()

			// container mount info
			for _, mount := range config.MountPoints {
				mountDetails = append(mountDetails, event.MountDetail{
					Destination: mount.Destination,
					Type:        mount.Type,
					Source:      mount.Source,
					Options:     []string{},
					Permission: func() string {
						if mount.Rw {
							return "rw"
						} else {
							return "ro"
						}
					}(),
					VolumeName: mount.Name,
				})
			}
		}
	case *containerd.Container:
		// skip moby namespace
		splits := strings.SplitN(c.ID(), "/", 2)
		if len(splits) == 2 {
			if splits[0] == "moby" {
				return nil
			}
		}

		// runtime type
		containerRuntime = event.Containerd

		// runtime desc
		runtimeUniqDesc = c.Runtime().UniqueDesc()

		// container mount info
		for _, mount := range ocispec.Mounts {
			permission := "rw"

			for _, option := range mount.Options {
				if option == "ro" {
					permission = "ro"
					break
				}
			}

			mountDetails = append(mountDetails, event.MountDetail{
				Destination: mount.Destination,
				Type:        mount.Type,
				Source:      mount.Source,
				Options:     mount.Options,
				VolumeName:  "-",
				Permission:  permission,
			})
		}
	}

	// root process
	if ocispec.Process != nil {
		rootProcessDetail.Terminal = ocispec.Process.Terminal
		rootProcessDetail.Env = ocispec.Process.Env
		rootProcessDetail.UID = ocispec.Process.User.UID
		rootProcessDetail.GID = ocispec.Process.User.GID
		rootProcessDetail.Args = ocispec.Process.Args
		rootProcessDetail.Cwd = ocispec.Process.Cwd

		if ocispec.Process.Capabilities != nil {
			rootProcessDetail.Capabilities = event.CapabilitiesDetail{
				Bounding:    ocispec.Process.Capabilities.Bounding,
				Effective:   ocispec.Process.Capabilities.Effective,
				Inheritable: ocispec.Process.Capabilities.Inheritable,
				Permitted:   ocispec.Process.Capabilities.Permitted,
				Ambient:     ocispec.Process.Capabilities.Ambient,
			}
		}

		// mapping username and groupname
		{
			entries, err := passwd.ParseFilesystemPasswd(container)
			if err != nil {
				log.Error(err)
			} else {
				for _, e := range entries {
					uid, err := strconv.ParseUint(e.Uid, 10, 32)
					if err != nil {
						log.Error(err)
						continue
					}

					if uint32(uid) == ocispec.Process.User.UID {
						rootProcessDetail.Username = e.Username
						break
					}
				}
			}
		}

		{
			entries, err := group.ParseFilesystemGroup(container)
			if err != nil {
				log.Error(err)
			} else {
				for _, e := range entries {
					gid, err := strconv.ParseUint(e.Gid, 10, 32)
					if err != nil {
						log.Error(err)
						continue
					}

					if uint32(gid) == ocispec.Process.User.GID {
						rootProcessDetail.Groupname = e.GroupName
						break
					}
				}
			}
		}
	}

	// container process
	pids, err := container.Pids()
	if err != nil {
		log.Error(err)
	} else {
		for _, pid := range pids {
			p, err := container.NewProcess(pid)
			if err != nil {
				log.Error(err)
				continue
			}

			cmdline, _ := p.Cmdline()
			cwd, _ := p.Cwd()
			env, _ := p.Environ()
			exe, _ := p.Exe()
			gids, _ := p.Gids()
			uids, _ := p.Uids()
			ppid, _ := p.Ppid()
			nspid, _ := p.Pid()
			hostPid, _ := p.HostPid()
			name, _ := p.Name()
			status, _ := p.Status()
			createTime, _ := p.CreateTime()
			p.Close()

			// mapping username and groupname
			usernames := make([]string, 4)
			{
				entries, err := passwd.ParseFilesystemPasswd(container)
				if err != nil {
					log.Error(err)
				} else {
					for _, e := range entries {
						uid, err := strconv.ParseInt(e.Uid, 10, 32)
						if err != nil {
							log.Error(err)
							continue
						}

						for index, uidT := range uids {
							if int32(uid) == uidT {
								usernames[index] = e.Username
							}
						}
					}
				}
			}

			groupnames := make([]string, 4)
			{
				entries, err := group.ParseFilesystemGroup(container)
				if err != nil {
					log.Error(err)
				} else {
					for _, e := range entries {
						gid, err := strconv.ParseInt(e.Gid, 10, 32)
						if err != nil {
							log.Error(err)
							continue
						}

						for index, gidT := range gids {
							if int32(gid) == gidT {
								groupnames[index] = e.GroupName
							}
						}
					}
				}
			}

			processDetails = append(processDetails, event.ProcessDetail{
				Cmdline:    cmdline,
				Cwd:        cwd,
				Environ:    env,
				Exe:        exe,
				Gids:       gids,
				Groupnames: groupnames,
				Uids:       uids,
				Usernames:  usernames,
				Pid:        nspid,
				Ppid:       ppid,
				HostPid:    hostPid,
				Status:     status,
				Name:       name,
				CreateTime: createTime.Unix(),
			})
		}
	}

	evt := &event.Event{
		BasicInfo: &event.BasicInfo{
			ID:         container.ID(),
			Object:     event.NewObject(container),
			Time:       time.Now(),
			Level:      event.Low,
			DetectType: event.Container,
			AlertType:  event.BasicContainer,
			EventType:  event.Info,
		},
		DetailInfo: &event.DetailInfo{
			AlertDetail: &event.ContainerBasicDetail{
				Name:            container.Name(),
				CreatedTime:     createdTime,
				State:           string(ocistate.Status),
				Runtime:         containerRuntime,
				RuntimeUniqDesc: runtimeUniqDesc,
				Hostname:        ocispec.Hostname,
				ImageID:         container.ImageID(),
				Privileged:      capability.IsPrivileged(container),
				RootProcess:     rootProcessDetail,
				Mounts:          mountDetails,
				Processes:       processDetails,
			},
		},
	}

	err = reportService.Client.Report(evt)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCommand.AddCommand(scanCommand)
	rootCommand.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-basic",
		Author:      "veinmind-team",
		Description: "veinmind-basic scan image basic info",
	}))
	scanCommand.AddCommand(report.MapReportCmd(cmd.MapImageCommand(scanImageCommand, scanImage), reportService, service.WithVerbose()))
	scanCommand.AddCommand(report.MapReportCmd(cmd.MapContainerCommand(scanContainerCommand, scanContainer), reportService, service.WithVerbose()))
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
