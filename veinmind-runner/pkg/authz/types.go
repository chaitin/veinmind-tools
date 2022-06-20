package authz

var ActionCheckType = make(map[string]string)

var (
	// ActionContainerArchive describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#get-an-archive-of-a-filesystem-resource-in-a-container
	ActionContainerArchive = "container_archive"
	// ActionContainerArchiveExtract describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#extract-an-archive-of-files-or-folders-to-a-directory-in-a-container
	ActionContainerArchiveExtract = "container_archive_extract"
	// ActionContainerArchiveInfo describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#retrieving-information-about-files-and-folders-in-a-container
	ActionContainerArchiveInfo = "container_archive_info"
	// ActionContainerAttach describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#attach-to-a-container
	ActionContainerAttach = "container_attach"
	// ActionContainerAttachWs describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#attach-to-a-container-websocket
	ActionContainerAttachWs = "container_attach_websocket"
	// ActionContainerChanges describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#inspect-changes-on-a-container-s-filesystem
	ActionContainerChanges = "container_changes"
	// ActionContainerCommit describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-a-new-image-from-a-container-s-changes
	ActionContainerCommit = "container_commit"
	// ActionContainerCopyFiles describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#copy-files-or-folders-from-a-container
	ActionContainerCopyFiles = "container_copyfiles"
	// ActionContainerCreate describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-a-container
	ActionContainerCreate = "container_create"
	// ActionContainerDelete describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#remove-a-container
	ActionContainerDelete = "container_delete"
	// ActionContainerExecCreate describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#exec-create
	ActionContainerExecCreate = "container_exec_create"
	// ActionContainerExecInspect describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#exec-inspect
	ActionContainerExecInspect = "container_exec_inspect"
	// ActionContainerExecStart describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#exec-start
	ActionContainerExecStart = "container_exec_start"
	// ActionContainerExport describes http://docs.docker.com/reference/api/docker_remote_api_v1.21/#export-a-container
	ActionContainerExport = "container_export"
	// ActionContainerInspect describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#inspect-a-container
	ActionContainerInspect = "container_inspect"
	// ActionContainerKill describes http://docs.docker.com/reference/api/docker_remote_api_v1.21/#kill-a-container
	ActionContainerKill = "container_kill"
	// ActionContainerList describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#list-containers
	ActionContainerList = "container_list"
	// ActionContainerLogs describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#get-container-logs
	ActionContainerLogs = "container_logs"
	// ActionContainerPause describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#pause-a-container
	ActionContainerPause = "container_pause"
	// ActionContainerRename describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#rename-a-container
	ActionContainerRename = "container_rename"
	// ActionContainerResize describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#resize-a-container-tty
	ActionContainerResize = "container_resize"
	// ActionContainerRestart describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#restart-a-container
	ActionContainerRestart = "container_restart"
	// ActionContainerStart describes http://docs.docker.com/reference/api/docker_remote_api_v1.21/#start-a-container
	ActionContainerStart = "container_start"
	// ActionContainerStats describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#get-container-stats-based-on-resource-usage
	ActionContainerStats = "container_stats"
	// ActionContainerStop describes http://docs.docker.com/reference/api/docker_remote_api_v1.21/#export-a-container
	ActionContainerStop = "container_stop"
	// ActionContainerTop describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#list-processes-running-inside-a-container
	ActionContainerTop = "container_top"
	// ActionContainerUnpause describes http://docs.docker.com/reference/api/docker_remote_api_v1.21/#unpause-a-container
	ActionContainerUnpause = "container_unpause"
	// ActionContainerWait describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#wait-a-container
	ActionContainerWait = "container_wait"
	// ActionDockerCheckAuth describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#check-auth-configuration
	ActionDockerCheckAuth = "docker_auth"
	// ActionDockerEvents describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#monitor-docker-s-events
	ActionDockerEvents = "docker_events"
	// ActionDockerInfo describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#display-system-wide-information
	ActionDockerInfo = "docker_info"
	// ActionDockerPing describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#ping-the-docker-server
	ActionDockerPing = "docker_ping"
	// ActionDockerVersion describes https://docs.docker.com/reference/api/docker_remote_api_v1.20/#show-the-docker-version-information
	ActionDockerVersion = "docker_version"
	// ActionImageArchive describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#get-a-tarball-containing-all-images
	ActionImageArchive = "images_archive"
	// ActionImageBuild describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#build-image-from-a-dockerfile
	ActionImageBuild = "image_build"
	// ActionImageCreate describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-an-image
	ActionImageCreate = "image_create"
	// ActionImageDelete describes https://docs.docker.com/reference/api/docker_remote_api_v1.18/#inspect-an-image
	ActionImageDelete = "image_delete"
	// ActionImageHistory describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#get-the-history-of-an-image
	ActionImageHistory = "image_history"
	// ActionImageInspect describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#inspect-an-image
	ActionImageInspect = "image_inspect"
	// ActionImageList describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#list-images
	ActionImageList = "image_list"
	// ActionImageLoad describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#load-a-tarball-with-a-set-of-images-and-tags-into-docker
	ActionImageLoad = "images_load"
	// ActionImagePrune describes https://docs.docker.com/engine/api/v1.37/#operation/ImagePrune
	ActionImagePrune = "image_prune"
	// ActionImagePush describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#push-an-image-on-the-registry
	ActionImagePush = "image_push"
	// ActionImagesSearch describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#search-images
	ActionImagesSearch = "images_search"
	// ActionImageTag describes https://docs.docker.com/reference/api/docker_remote_api_v1.21/#tag-an-image-into-a-repository
	ActionImageTag = "image_tag"
	// ActionVolumeList describes  https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#list-volumes
	ActionVolumeList = "volume_list"
	// ActionVolumeCreate describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-a-volume
	ActionVolumeCreate = "volume_create"
	// ActionVolumeInspect describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#inspect-a-volume
	ActionVolumeInspect = "volume_inspect"
	// ActionVolumeRemove describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#remove-a-volume
	ActionVolumeRemove = "volume_remove"
	// ActionNetworkList describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#list-networks
	ActionNetworkList = "network_list"
	// ActionNetworkInspect describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#inspect-network
	ActionNetworkInspect = "network_inspect"
	// ActionNetworkCreate describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-a-network
	ActionNetworkCreate = "network_create"
	// ActionNetworkConnect describes
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#connect-a-container-to-a-network
	ActionNetworkConnect = "network_connect"
	// ActionNetworkDisconnect describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#disconnect-a-container-from-a-network
	ActionNetworkDisconnect = "network_disconnect"
	// ActionNetworkRemove describes https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#remove-a-network
	ActionNetworkRemove = "network_remove"
	// ActionSwarmInspect describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmInspect
	ActionSwarmInspect = "swarm_inspect"
	// ActionSwarmInit describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmInit
	ActionSwarmInit = "swarm_init"
	// ActionSwarmJoin describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmJoin
	ActionSwarmJoin = "swarm_join"
	// ActionSwarmLeave describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmLeave
	ActionSwarmLeave = "swarm_leave"
	// ActionSwarmUpdate describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmUpdate
	ActionSwarmUpdate = "swarm_update"
	// ActionSwarmUnlockKey describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmUnlockkey
	ActionSwarmUnlockKey = "swarm_unlock_key"
	// ActionSwarmUnlock describes https://docs.docker.com/engine/api/v1.37/#operation/SwarmUnlock
	ActionSwarmUnlock = "swarm_unlock"
	// ActionNodeList describes https://docs.docker.com/engine/api/v1.39/#operation/NodeList
	ActionNodeList = "node_list"
	// ActionNodeInspect describes https://docs.docker.com/engine/api/v1.39/#operation/NodeInspect
	ActionNodeInspect = "node_inspect"
	// ActionNodeDelete describes https://docs.docker.com/engine/api/v1.39/#operation/NodeDelete
	ActionNodeDelete = "node_delete"
	// ActionNodeUpdate describes https://docs.docker.com/engine/api/v1.39/#operation/NodeUpdate
	ActionNodeUpdate = "node_update"
	// ActionServiceList describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceList
	ActionServiceList = "service_list"
	// ActionServiceCreate describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceCreate
	ActionServiceCreate = "service_create"
	// ActionServiceInspect describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceInspect
	ActionServiceInspect = "service_inspect"
	// ActionServiceDelete describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceDelete
	ActionServiceDelete = "service_delete"
	// ActionServiceUpdate describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceUpdate
	ActionServiceUpdate = "service_update"
	// ActionServiceLogs describes https://docs.docker.com/engine/api/v1.39/#operation/ServiceLogs
	ActionServiceLogs = "service_logs"
	// ActionTaskList describes https://docs.docker.com/engine/api/v1.39/#operation/TaskList
	ActionTaskList = "task_list"
	// ActionTaskInspect describes https://docs.docker.com/engine/api/v1.39/#operation/TaskInspect
	ActionTaskInspect = "task_inspect"
	// ActionSecretList describes https://docs.docker.com/engine/api/v1.39/#operation/SecretList
	ActionSecretList = "secret_list"
	// ActionSecretCreate describes https://docs.docker.com/engine/api/v1.39/#operation/SecretCreate
	ActionSecretCreate = "secret_create"
	// ActionSecretInspect describes https://docs.docker.com/engine/api/v1.39/#operation/SecretInspect
	ActionSecretInspect = "secret_inspect"
	// ActionSecretDelete describes https://docs.docker.com/engine/api/v1.39/#operation/SecretDelete
	ActionSecretDelete = "secret_delete"
	// ActionSecretUpdate describes https://docs.docker.com/engine/api/v1.39/#operation/SecretUpdate
	ActionSecretUpdate = "secret_update"
	// ActionConfigList describes https://docs.docker.com/engine/api/v1.39/#operation/ConfigList
	ActionConfigList = "config_list"
	// ActionConfigCreate describes https://docs.docker.com/engine/api/v1.39/#operation/ConfigCreate
	ActionConfigCreate = "config_create"
	// ActionConfigInspect describes https://docs.docker.com/engine/api/v1.39/#operation/ConfigInspect
	ActionConfigInspect = "config_inspect"
	// ActionConfigDelete describes https://docs.docker.com/engine/api/v1.39/#operation/ConfigDelete
	ActionConfigDelete = "config_delete"
	// ActionConfigUpdate describes https://docs.docker.com/engine/api/v1.39/#operation/ConfigUpdate
	ActionConfigUpdate = "config_update"
	// ActionDistributionInspect describes https://docs.docker.com/engine/api/v1.39/#operation/DistributionInspect
	ActionDistributionInspect = "distribution_inspect"
	// ActionNone indicates no action matched the given method URL combination
	ActionNone = ""
)
