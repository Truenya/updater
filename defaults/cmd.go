package defaults

import "fmt"

const DefaultRemoteDir = "/tmp/update_dir/"
const DefaultDir = "/opt/ispsystem/"
const ContainerDeleteCmd = "docker rm %s" // container_name

func CreateBinDirCmd(service string) string {
	return CreateDirCmd(service + "/bin")
}

func CreateDirCmd(service string) string {
	return fmt.Sprintf("mkdir -p %s%s", DefaultRemoteDir, service)
}

func RemoveExistingDirCmd(service string) string {
	return fmt.Sprintf("rm -rf %s%s/*", DefaultRemoteDir, service)
}

func PullImageCmd(service string) string { // branch
	return fmt.Sprintf("docker pull registry-dev.ispsystem.net/team/vm/%s", service+":%s")
}

func CreateContainerCmd(service string) string { // container_name, branch
	return fmt.Sprintf("docker create --name %s registry-dev.ispsystem.net/team/vm/%s sh", "%s", service+":%s")
}

func RsyncArg(second string) []string {
	return []string{"", second}
}

func RsyncArgs(first, second string) []string {
	return []string{first, second}
}

func DockerCpBinToBoxDefaultCmd(service string) string {
	return DockerCpBinToDefaultCmd(service, "vm_box")
}

func DockerCpBinToDefaultCmd(service, container string) string {
	return DockerCpToDefaultCmd(service, "bin/", container)
}

func DockerCpToDefaultBoxCmd(service, what string) string {
	return DockerCpToDefaultCmd(service, what, "vm_box")
}

func DockerCpToDefaultCmd(service, what, container string) string {
	return DockerCpToDefaultDestCmd(service, what, container, service)
}

func DockerCpToBoxDestCmd(service, what, to string) string {
	return DockerCpToDestCmd(service, what, "vm_box", to)
}

func DockerCpToDefaultBoxDestCmd(service, what, to string) string {
	return DockerCpToDefaultDestCmd(service, what, "vm_box", to)
}

func DockerCpToDefaultDestCmd(service, what, container, to string) string {
	return fmt.Sprintf("docker cp %s%s/%s %s:%s%s", DefaultRemoteDir, service, what, container, DefaultDir, to)
}

func DockerCpToDestCmd(service, what, container, to string) string {
	return fmt.Sprintf("docker cp %s%s/%s %s:%s", DefaultRemoteDir, service, what, container, to)
}

func DockerCpToCmd(from, to, container string) string {
	return fmt.Sprintf("docker cp %s %s:%s", from, container, to)
}

func DockerRestartBoxServiceCmd(service string) string {
	return DockerRestartServiceCmd(service, "vm_box")
}

func DockerRestartServiceCmd(service, container string) string {
	return fmt.Sprintf("docker exec %s supervisorctl restart %s", container, service)
}

func DockerCpBinFromCmd(service, what string) string { // container_name
	return fmt.Sprintf("docker cp %s %s%s/bin/", "%s:/"+what, DefaultRemoteDir, service)
}

func DockerCpServiceFromCmd(service, what string) string { // container_name
	return DockerCpFromToCmd(DefaultDir+what, service) // what - папка в которой сервис лежит в дефолтной дире
}

func DockerCpFromCmd(from string) string { // container_name
	return fmt.Sprintf("docker cp %s %s", "%s:"+from, DefaultRemoteDir)
}

func DockerCpFromToCmd(from, to string) string { // container_name
	return fmt.Sprintf("docker cp %s %s%s", "%s:"+from, DefaultRemoteDir, to)
}

func MoveServiceDefDir(service, from, to string) string {
	return fmt.Sprintf("mv %s%s/%s %s%s/%s", DefaultRemoteDir, service, from, DefaultRemoteDir, service, to)
}
