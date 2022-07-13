package action

type DockerPlugin string

var (
	//  ContainerArchive describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#get-an-archive-of-a-filesystem-resource-in-a-container
	ContainerArchive DockerPlugin = "container_archive"
	//  ContainerArchiveExtract describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#extract-an-archive-of-files-or-folders-to-a-directory-in-a-container
	ContainerArchiveExtract DockerPlugin = "container_archive_extract"
	//  ContainerArchiveInfo describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#retrieving-information-about-files-and-folders-in-a-container
	ContainerArchiveInfo DockerPlugin = "container_archive_info"
	//  ContainerAttach describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#attach-to-a-container
	ContainerAttach DockerPlugin = "container_attach"
	//  ContainerAttachWs describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#attach-to-a-container-websocket
	ContainerAttachWs DockerPlugin = "container_attach_websocket"
	//  ContainerChanges describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#inspect-changes-on-a-container-s-filesystem
	ContainerChanges DockerPlugin = "container_changes"
	//  ContainerCommit describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-a-new-image-from-a-container-s-changes
	ContainerCommit DockerPlugin = "container_commit"
	//  ContainerCopyFiles describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#copy-files-or-folders-from-a-container
	ContainerCopyFiles DockerPlugin = "container_copyfiles"
	//  ContainerCreate describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-a-container
	ContainerCreate DockerPlugin = "container_create"
	//  ContainerDelete describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#remove-a-container
	ContainerDelete DockerPlugin = "container_delete"
	//  ContainerExecCreate describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#exec-create
	ContainerExecCreate DockerPlugin = "container_exec_create"
	//  ContainerExecInspect describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#exec-inspect
	ContainerExecInspect DockerPlugin = "container_exec_inspect"
	//  ContainerExecStart describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#exec-start
	ContainerExecStart DockerPlugin = "container_exec_start"
	//  ContainerExport describes http://docs.docker.com/reference/api/docker_remote_api_v1.21/#export-a-container
	ContainerExport DockerPlugin = "container_export"
	//  ContainerInspect describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#inspect-a-container
	ContainerInspect DockerPlugin = "container_inspect"
	//  ContainerKill describes http://docs.docker.com/reference/api/docker_remote_api_v1.21/#kill-a-container
	ContainerKill DockerPlugin = "container_kill"
	//  ContainerList describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#list-containers
	ContainerList DockerPlugin = "container_list"
	//  ContainerLogs describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#get-container-logs
	ContainerLogs DockerPlugin = "container_logs"
	//  ContainerPause describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#pause-a-container
	ContainerPause DockerPlugin = "container_pause"
	//  ContainerRename describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#rename-a-container
	ContainerRename DockerPlugin = "container_rename"
	//  ContainerResize describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#resize-a-container-tty
	ContainerResize DockerPlugin = "container_resize"
	//  ContainerRestart describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#restart-a-container
	ContainerRestart DockerPlugin = "container_restart"
	//  ContainerStart describes http://docs.docker.com/reference/api/docker_remote_api_v1.21/#start-a-container
	ContainerStart DockerPlugin = "container_start"
	//  ContainerStats describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#get-container-stats-based-on-resource-usage
	ContainerStats DockerPlugin = "container_stats"
	//  ContainerStop describes http://docs.docker.com/reference/api/docker_remote_api_v1.21/#export-a-container
	ContainerStop DockerPlugin = "container_stop"
	//  ContainerTop describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#list-processes-running-inside-a-container
	ContainerTop DockerPlugin = "container_top"
	//  ContainerUnpause describes http://docs.docker.com/reference/api/docker_remote_api_v1.21/#unpause-a-container
	ContainerUnpause DockerPlugin = "container_unpause"
	//  ContainerWait describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#wait-a-container
	ContainerWait DockerPlugin = "container_wait"
	//  DockerCheckAuth describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#check-auth-configuration
	DockerCheckAuth DockerPlugin = "docker_auth"
	//  DockerEvents describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#monitor-docker-s-events
	DockerEvents DockerPlugin = "docker_events"
	//  DockerInfo describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#display-system-wide-information
	DockerInfo DockerPlugin = "docker_info"
	//  DockerPing describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#ping-the-docker-server
	DockerPing DockerPlugin = "docker_ping"
	//  DockerVersion describes https://docs.docker.com/reference/api/docker_remote_api_v1.20/#show-the-docker-version-information
	DockerVersion DockerPlugin = "docker_version"
	//  ImageArchive describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#get-a-tarball-containing-all-images
	ImageArchive DockerPlugin = "images_archive"
	//  ImageBuild describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#build-image-from-a-dockerfile
	ImageBuild DockerPlugin = "image_build"
	//  ImageCreate describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-an-image
	ImageCreate DockerPlugin = "image_create"
	//  ImageDelete describes https://docs.docker.com/reference/api/docker_remote_api_v1.18/#inspect-an-image
	ImageDelete DockerPlugin = "image_delete"
	//  ImageHistory describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#get-the-history-of-an-image
	ImageHistory DockerPlugin = "image_history"
	//  ImageInspect describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#inspect-an-image
	ImageInspect DockerPlugin = "image_inspect"
	//  ImageList describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#list-images
	ImageList DockerPlugin = "image_list"
	//  ImageLoad describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#load-a-tarball-with-a-set-of-images-and-tags-into-docker
	ImageLoad DockerPlugin = "images_load"
	//  ImagePrune describes https://docs.docker.com/engine/api/v1.37/#operation/ImagePrune
	ImagePrune DockerPlugin = "image_prune"
	//  ImagePush describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#push-an-image-on-the-registry
	ImagePush DockerPlugin = "image_push"
	//  ImagesSearch describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#search-images
	ImagesSearch DockerPlugin = "images_search"
	//  ImageTag describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#tag-an-image-into-a-repository
	ImageTag DockerPlugin = "image_tag"
	//  VolumeList describes  https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#list-volumes
	VolumeList DockerPlugin = "volume_list"
	//  VolumeCreate describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-a-volume
	VolumeCreate DockerPlugin = "volume_create"
	//  VolumeInspect describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#inspect-a-volume
	VolumeInspect DockerPlugin = "volume_inspect"
	//  VolumeRemove describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#remove-a-volume
	VolumeRemove DockerPlugin = "volume_remove"
	//  NetworkList describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#list-networks
	NetworkList DockerPlugin = "network_list"
	//  NetworkInspect describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#inspect-network
	NetworkInspect DockerPlugin = "network_inspect"
	//  NetworkCreate describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-a-network
	NetworkCreate DockerPlugin = "network_create"
	//  NetworkConnect describes
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#connect-a-container-to-a-network
	NetworkConnect DockerPlugin = "network_connect"
	//  NetworkDisconnect describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#disconnect-a-container-from-a-network
	NetworkDisconnect DockerPlugin = "network_disconnect"
	//  NetworkRemove describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#remove-a-network
	NetworkRemove DockerPlugin = "network_remove"
	//  SwarmInspect describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmInspect
	SwarmInspect DockerPlugin = "swarm_inspect"
	//  SwarmInit describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmInit
	SwarmInit DockerPlugin = "swarm_init"
	//  SwarmJoin describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmJoin
	SwarmJoin DockerPlugin = "swarm_join"
	//  SwarmLeave describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmLeave
	SwarmLeave DockerPlugin = "swarm_leave"
	//  SwarmUpdate describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmUpdate
	SwarmUpdate DockerPlugin = "swarm_update"
	//  SwarmUnlockKey describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmUnlockkey
	SwarmUnlockKey DockerPlugin = "swarm_unlock_key"
	//  SwarmUnlock describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmUnlock
	SwarmUnlock DockerPlugin = "swarm_unlock"
	//  NodeList describes https://docs.docker.com/engine/api/v1.39/#operation/NodeList
	NodeList DockerPlugin = "node_list"
	//  NodeInspect describes https://docs.docker.com/engine/api/v1.39/#operation/NodeInspect
	NodeInspect DockerPlugin = "node_inspect"
	//  NodeDelete describes https://docs.docker.com/engine/api/v1.39/#operation/NodeDelete
	NodeDelete DockerPlugin = "node_delete"
	//  NodeUpdate describes https://docs.docker.com/engine/api/v1.39/#operation/NodeUpdate
	NodeUpdate DockerPlugin = "node_update"
	//  ServiceList describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceList
	ServiceList DockerPlugin = "service_list"
	//  ServiceCreate describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceCreate
	ServiceCreate DockerPlugin = "service_create"
	//  ServiceInspect describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceInspect
	ServiceInspect DockerPlugin = "service_inspect"
	//  ServiceDelete describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceDelete
	ServiceDelete DockerPlugin = "service_delete"
	//  ServiceUpdate describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceUpdate
	ServiceUpdate DockerPlugin = "service_update"
	//  ServiceLogs describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceLogs
	ServiceLogs DockerPlugin = "service_logs"
	//  TaskList describes https://docs.docker.com/engine/api/v1.39/#operation/TaskList
	TaskList DockerPlugin = "task_list"
	//  TaskInspect describes https://docs.docker.com/engine/api/v1.39/#operation/TaskInspect
	TaskInspect DockerPlugin = "task_inspect"
	//  SecretList describes https://docs.docker.com/engine/api/v1.39/#operation/SecretList
	SecretList DockerPlugin = "secret_list"
	//  SecretCreate describes https://docs.docker.com/engine/api/v1.39/#operation/SecretCreate
	SecretCreate DockerPlugin = "secret_create"
	//  SecretInspect describes https://docs.docker.com/engine/api/v1.39/#operation/SecretInspect
	SecretInspect DockerPlugin = "secret_inspect"
	//  SecretDelete describes https://docs.docker.com/engine/api/v1.39/#operation/SecretDelete
	SecretDelete DockerPlugin = "secret_delete"
	//  SecretUpdate describes https://docs.docker.com/engine/api/v1.39/#operation/SecretUpdate
	SecretUpdate DockerPlugin = "secret_update"
	//  ConfigList describes https://docs.docker.com/engine/api/v1.39/#operation/ConfigList
	ConfigList DockerPlugin = "config_list"
	//  ConfigCreate describes https://docs.docker.com/engine/api/v1.39/#operation/ConfigCreate
	ConfigCreate DockerPlugin = "config_create"
	//  ConfigInspect describes https://docs.docker.com/engine/api/v1.39/#operation/ConfigInspect
	ConfigInspect DockerPlugin = "config_inspect"
	//  ConfigDelete describes https://docs.docker.com/engine/api/v1.39/#operation/ConfigDelete
	ConfigDelete DockerPlugin = "config_delete"
	//  ConfigUpdate describes https://docs.docker.com/engine/api/v1.39/#operation/ConfigUpdate
	ConfigUpdate DockerPlugin = "config_update"
	//  DistributionInspect describes https://docs.docker.com/engine/api/v1.39/#operation/DistributionInspect
	DistributionInspect DockerPlugin = "distribution_inspect"
	//  None indicates no routeaction matched the given method URL combination
	NoneDockerPlugin DockerPlugin = ""
)
