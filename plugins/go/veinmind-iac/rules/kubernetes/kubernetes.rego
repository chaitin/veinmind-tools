package brightMirror.kubernetes

# kubernetes lib

pods[pod] {
	is_pod
	pod = input[_]
}

pods[pod] {
	is_controller
	pod = input[_].spec.template
}

pods[pod] {
	is_cronjob
	pod = input[_].spec.jobTemplate.spec.template
}

containers[container] {
	is_pod
    keys = {"containers", "initContainers"}
    all_containers = [c | keys[k]; c = input[_].spec[k][_]]

	container = all_containers[_]
}

volumes[volume] {
    is_pod
    volume = input[_].spec.volumes[_]
}

annotations[annotation] {
    pods[pod]
	annotation := pod.metadata.annotations
}

# container_securityContext
securityContexts[sec] {
    sec := containers[_].securityContext
}

# spec_securityContext
securityContexts[sec] {
    sec := pods[_].spec.securityContext
}

allowPrivilegeEscalations[allow] {
    allow := securityContexts[_].allowPrivilegeEscalations
}

is_pod {
    input[_].kind == "Pod"
}

is_cronjob {
	input[_].kind = "CronJob"
}

default is_controller = false

is_controller {
	input[_].kind = "Deployment"
}

is_controller {
	input[_].kind = "StatefulSet"
}

is_controller {
	input[_].kind = "DaemonSet"
}

is_controller {
	input[_].kind = "ReplicaSet"
}

is_controller {
	input[_].kind = "ReplicationController"
}

is_controller {
	input[_].kind = "Job"
}
