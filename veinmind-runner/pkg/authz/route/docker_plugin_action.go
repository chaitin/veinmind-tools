package route

import (
	"regexp"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/action"
	"github.com/docker/docker/pkg/authorization"
)

type routeAction struct {
	pattern string
	method  string
	action  action.DockerPlugin
}

var routes = []routeAction{
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#build-image-from-a-dockerfile
	{pattern: "/build", method: "POST", action: action.ImageBuild},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.20/#create-a-new-image-from-a-container-s-changes
	{pattern: "/commit", method: "POST", action: action.ContainerCommit},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.20/#monitor-docker-s-events
	{pattern: "/events", method: "POST", action: action.DockerEvents},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.20/#show-the-docker-version-information
	{pattern: "/version", method: "GET", action: action.DockerVersion},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.20/#check-auth-configuration
	{pattern: "/auth", method: "POST", action: action.DockerCheckAuth},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#wait-a-container
	{pattern: "/containers/.+/wait", method: "POST", action: action.ContainerWait},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#resize-a-container-tty
	{pattern: "/containers/.+/resize", method: "POST", action: action.ContainerResize},
	// http://docs.docker.com/reference/api/docker_remote_api_v1.21/#export-a-container
	{pattern: "/containers/.+/export", method: "POST", action: action.ContainerExport},
	// http://docs.docker.com/reference/api/docker_remote_api_v1.21/#export-a-container
	{pattern: "/containers/.+/stop", method: "POST", action: action.ContainerStop},
	// http://docs.docker.com/reference/api/docker_remote_api_v1.21/#kill-a-container
	{pattern: "/containers/.*/kill", method: "POST", action: action.ContainerKill},
	// http://docs.docker.com/reference/api/docker_remote_api_v1.21/#restart-a-container
	{pattern: "/containers/.+/restart", method: "POST", action: action.ContainerRestart},
	// http://docs.docker.com/reference/api/docker_remote_api_v1.21/#start-a-container
	{pattern: "/containers/.+/start", method: "POST", action: action.ContainerStart},
	// http://docs.docker.com/reference/api/docker_remote_api_v1.21/#exec-create
	{pattern: "/containers/.+/exec", method: "POST", action: action.ContainerExecCreate},
	// http://docs.docker.com/reference/api/docker_remote_api_v1.21/#unpause-a-container
	{pattern: "/containers/.+/unpause", method: "POST", action: action.ContainerUnpause},
	// http://docs.docker.com/reference/api/docker_remote_api_v1.21/#pause-a-container
	{pattern: "/containers/.+/pause", method: "POST", action: action.ContainerPause},
	// http://docs.docker.com/reference/api/docker_remote_api_v1.21/#copy-files-or-folders-from-a-container
	{pattern: "/containers/.+/copy", method: "POST", action: action.ContainerCopyFiles},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#extract-an-archive-of-files-or-folders-to-a-directory-in-a-container
	{pattern: "/containers/.+/archive", method: "PUT", action: action.ContainerArchiveExtract},
	{pattern: "/containers/.+/archive", method: "HEAD", action: action.ContainerArchiveInfo},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#get-an-archive-of-a-filesystem-resource-in-a-container
	{pattern: "/containers/.+/archive", method: "GET", action: action.ContainerArchive},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#attach-to-a-container-websocket
	{pattern: "/containers/.+/attach/ws", method: "GET", action: action.ContainerAttachWs},
	// http://docs.docker.com/reference/api/docker_remote_api_v1.21/#attach-to-a-container
	{pattern: "/containers/.+/attach", method: "POST", action: action.ContainerAttach},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#list-containers
	{pattern: "/containers/json", method: "GET", action: action.ContainerList},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#inspect-a-container
	{pattern: "/containers/.+/json", method: "GET", action: action.ContainerInspect},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#remove-a-container
	{pattern: "/containers/.+", method: "DELETE", action: action.ContainerDelete},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#rename-a-container
	{pattern: "/containers/.+/rename", method: "POST", action: action.ContainerRename},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#get-container-stats-based-on-resource-usage
	{pattern: "/containers/.+/stats", method: "GET", action: action.ContainerStats},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#inspect-changes-on-a-container-s-filesystem
	{pattern: "/containers/.+/changes", method: "GET", action: action.ContainerChanges},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#list-processes-running-inside-a-container
	{pattern: "/containers/.+/top", method: "GET", action: action.ContainerTop},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#get-container-logs
	{pattern: "/containers/.+/logs", method: "GET", action: action.ContainerLogs},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#create-a-container
	{pattern: "/containers/create", method: "POST", action: action.ContainerCreate},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#get-a-tarball-containing-all-images
	{pattern: "/images/.+./get", method: "GET", action: action.ImageArchive},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#search-images
	{pattern: "/images/search", method: "GET", action: action.ImagesSearch},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#tag-an-image-into-a-repository
	{pattern: "/images/.+/tag", method: "POST", action: action.ImageTag},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#inspect-an-image
	{pattern: "/images/.+/json", method: "GET", action: action.ImageInspect},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.18/#inspect-an-image
	{pattern: "/images/.+", method: "DELETE", action: action.ImageDelete},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#get-the-history-of-an-image
	{pattern: "/images/.+/history", method: "GET", action: action.ImageHistory},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#push-an-image-on-the-registry
	{pattern: "/images/.+/push", method: "POST", action: action.ImagePush},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#create-an-image
	{pattern: "/images/create", method: "POST", action: action.ImageCreate},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#load-a-tarball-with-a-set-of-images-and-tags-into-docker
	{pattern: "/images/load", method: "POST", action: action.ImageLoad},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#list-images
	{pattern: "/images/json", method: "GET", action: action.ImageList},
	// https://docs.docker.com/engine/api/v1.37/#operation/ImagePrune
	{pattern: "/images/prune", method: "POST", action: action.ImagePrune},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#ping-the-docker-server
	{pattern: "/_ping", method: "GET", action: action.DockerPing},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#display-system-wide-information
	{pattern: "/info", method: "GET", action: action.DockerInfo},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#exec-inspect
	{pattern: "/exec/.+/json", method: "GET", action: action.ContainerExecInspect},
	// https://docs.docker.com/reference/api/docker_remote_api_v1.21/#exec-start
	{pattern: "/exec/.+/start", method: "POST", action: action.ContainerExecStart},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#inspect-a-volume
	{pattern: "/volumes/.+", method: "GET", action: action.VolumeInspect},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#list-volumes
	{pattern: "/volumes", method: "GET", action: action.VolumeList},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-a-volume
	{pattern: "/volumes/create", method: "POST", action: action.VolumeCreate},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#remove-a-volume
	{pattern: "/volumes/.+", method: "DELETE", action: action.VolumeRemove},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#inspect-network
	{pattern: "/networks/.+", method: "GET", action: action.NetworkInspect},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#list-networks
	{pattern: "/networks", method: "GET", action: action.NetworkList},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#create-a-network
	{pattern: "/networks/create", method: "POST", action: action.NetworkCreate},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#connect-a-container-to-a-network
	{pattern: "/networks/.+/connect", method: "POST", action: action.NetworkConnect},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#disconnect-a-container-from-a-network
	{pattern: "/networks/.+/disconnect", method: "POST", action: action.NetworkDisconnect},
	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.21/#remove-a-network
	{pattern: "/networks/.+", method: "DELETE", action: action.NetworkRemove},
	// https://docs.docker.com/engine/api/v1.37/#operation/SwarmInit
	{pattern: "/swarm/init", method: "POST", action: action.SwarmInit},
	// https://docs.docker.com/engine/api/v1.37/#operation/SwarmJoin
	{pattern: "/swarm/join", method: "POST", action: action.SwarmJoin},
	// https://docs.docker.com/engine/api/v1.37/#operation/SwarmLeave
	{pattern: "/swarm/leave", method: "POST", action: action.SwarmLeave},
	// https://docs.docker.com/engine/api/v1.37/#operation/SwarmUpdate
	{pattern: "/swarm/update", method: "POST", action: action.SwarmUpdate},
	// https://docs.docker.com/engine/api/v1.37/#operation/SwarmUnlockkey
	{pattern: "/swarm/unlockkey", method: "GET", action: action.SwarmUnlockKey},
	// https://docs.docker.com/engine/api/v1.37/#operation/SwarmUnlock
	{pattern: "/swarm/unlock", method: "POST", action: action.SwarmUnlock},
	// https://docs.docker.com/engine/api/v1.37/#operation/SwarmInspect
	{pattern: "/swarm", method: "GET", action: action.SwarmInspect},
	// https://docs.docker.com/engine/api/v1.39/#operation/NodeUpdate
	{pattern: "/nodes/.+/update", method: "POST", action: action.NodeUpdate},
	// https://docs.docker.com/engine/api/v1.39/#operation/NodeInspect
	{pattern: "/nodes/.+", method: "GET", action: action.NodeInspect},
	// https://docs.docker.com/engine/api/v1.39/#operation/NodeDelete
	{pattern: "/nodes/.+", method: "DELETE", action: action.NodeDelete},
	// https://docs.docker.com/engine/api/v1.39/#operation/NodeList
	{pattern: "/nodes", method: "GET", action: action.NodeList},
	// https://docs.docker.com/engine/api/v1.39/#operation/ServiceCreate
	{pattern: "/services/create", method: "POST", action: action.ServiceCreate},
	// https://docs.docker.com/engine/api/v1.39/#operation/ServiceUpdate
	{pattern: "/services/.+/update", method: "POST", action: action.ServiceUpdate},
	// https://docs.docker.com/engine/api/v1.39/#operation/ServiceLogs
	{pattern: "/services/.+/logs", method: "GET", action: action.ServiceLogs},
	// https://docs.docker.com/engine/api/v1.39/#operation/ServiceInspect
	{pattern: "/services/.+", method: "GET", action: action.ServiceInspect},
	// https://docs.docker.com/engine/api/v1.39/#operation/ServiceDelete
	{pattern: "/services/.+", method: "DELETE", action: action.ServiceDelete},
	// https://docs.docker.com/engine/api/v1.39/#operation/ServiceList
	{pattern: "/services", method: "GET", action: action.ServiceList},
	// https://docs.docker.com/engine/api/v1.39/#operation/TaskInspect
	{pattern: "/tasks/.+", method: "GET", action: action.TaskInspect},
	// https://docs.docker.com/engine/api/v1.39/#operation/TaskList
	{pattern: "/tasks", method: "GET", action: action.TaskList},
	// https://docs.docker.com/engine/api/v1.39/#operation/SecretCreate
	{pattern: "/secrets/create", method: "POST", action: action.SecretCreate},
	// https://docs.docker.com/engine/api/v1.39/#operation/SecretUpdate
	{pattern: "/secrets/.+/update", method: "POST", action: action.SecretUpdate},
	// https://docs.docker.com/engine/api/v1.39/#operation/SecretInspect
	{pattern: "/secrets/.+", method: "GET", action: action.SecretInspect},
	// https://docs.docker.com/engine/api/v1.39/#operation/SecretDelete
	{pattern: "/secrets/.+", method: "DELETE", action: action.SecretDelete},
	// https://docs.docker.com/engine/api/v1.39/#operation/SecretList
	{pattern: "/secrets", method: "GET", action: action.SecretList},
	// https://docs.docker.com/engine/api/v1.39/#operation/ConfigCreate
	{pattern: "/configs/create", method: "POST", action: action.ConfigCreate},
	// https://docs.docker.com/engine/api/v1.39/#operation/ConfigUpdate
	{pattern: "/configs/.+/update", method: "POST", action: action.ConfigUpdate},
	// https://docs.docker.com/engine/api/v1.39/#operation/ConfigInspect
	{pattern: "/configs/.+", method: "GET", action: action.ConfigInspect},
	// https://docs.docker.com/engine/api/v1.39/#operation/ConfigDelete
	{pattern: "/configs/.+", method: "DELETE", action: action.ConfigDelete},
	// https://docs.docker.com/engine/api/v1.39/#operation/ConfigList
	{pattern: "/configs", method: "GET", action: action.ConfigList},
	// https://docs.docker.com/engine/api/v1.39/#operation/DistributionInspect
	{pattern: "/distribution/.+/json", method: "GET", action: action.DistributionInspect},
}

var routeRegexes map[*regexp.Regexp]routeAction

func init() {
	routeRegexes = make(map[*regexp.Regexp]routeAction)
	for _, route := range routes {
		re := regexp.MustCompile(route.pattern)
		routeRegexes[re] = route
	}
}

// ParseRoute convert a method/url pattern to corresponding docker routeaction
func ParseDockerPluginAction(req *authorization.Request) action.DockerPlugin {
	for re, route := range routeRegexes {
		if route.method == req.RequestMethod {
			if re.MatchString(req.RequestURI) {
				return route.action
			}
		}
	}

	return action.NoneDockerPlugin
}
