package client

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
    "encoding/base64"
	"bufio"
    "bytes"

	"io"
    

    // Docker Engine API
	dockerapi "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/filters"
    
    // Consul API
    consulapi "github.com/hashicorp/consul/api"
 

	"github.com/samalba/dockerclient"

	ermtypes "releases-manager/orchestrator/types"
)

const (
	// Shell commands (in container)
	EXOR_CMD_RELEASE_START   string = "release-start"
	EXOR_CMD_DISPLAY_CATALOG string = "catalog-from-url"
)

// DockerClient the docker cient
var cli *dockerapi.Client
var dockererr, consulerr error
var consulCli *consulapi.Client

func init(){
    cli, dockererr = dockerapi.NewEnvClient()
    
    // Get a new Consul client
    config := &consulapi.Config{
		Address:    "192.168.99.100:8500",
		Scheme:     "http",
	}
    consulCli, consulerr = consulapi.NewClient(config)
    if consulerr != nil {
        panic(consulerr)
    }
}

func ExecAllReleases(baseCatalog ermtypes.Catalog, label string) {

	fmt.Println("start release-all")

	// filter projects by label and sort by container step
	var projectsToRelease ermtypes.Catalog

	for _, p := range baseCatalog {
		if strings.Contains(p.Labels, label) {
			projectsToRelease = append(projectsToRelease, p)
		}
	}
	log.Printf("p to release: %d", projectsToRelease.Len())

	sort.Sort(projectsToRelease)

	for _, p := range projectsToRelease {
		fmt.Println(p.Name)
		log.Printf("Started container for:%s - step: %d", p.Name, p.Container.Step)
		rs, err := ReadReleaseStatusFromContainer(p.Name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Name %s", rs.Name)
		//startReleaseContainer(cli, p.Name, p.Container.Image, cmd)
	}

}

// Create a volume for the workspace
func createVolumeWorkspace(name string) {

	var isVolumeAlreadyCreated = false
	var volumeName = "workspace-" + name

	var args = filters.NewArgs()
	args.Add("name", volumeName)
	list, err := cli.VolumeList(args)
	if err != nil {
		panic(err)
	}
	for _, v := range list.Volumes {
		if v.Name == volumeName {
			log.Print("Volume already exists.")
			isVolumeAlreadyCreated = true
		}

	}

	if isVolumeAlreadyCreated == false {
		log.Print("Create Workspace Volume.")
		var conf = types.VolumeCreateRequest{
			Name: "workspace-" + name,
		}
		cli.VolumeCreate(conf)
	}

}

// Start a eXo Release Container and exec a command
func startReleaseContainer(name string, image string, cmd string) {

	var containerName = "exor-" + name

	var hostConf = container.HostConfig{
		Binds: BindFiles(name),
	}
	var conf = container.Config{
		Env:   ConfigFileAsEnv(),
		Image: image,
		User:  EXOR_USER,
		Cmd:   []string{cmd},
	}

	c, err := cli.ContainerCreate(&conf, &hostConf, nil, containerName)
	if err != nil {
		panic(err)
	}
	err2 := cli.ContainerStart(c.ID)
	if err2 != nil {
		panic(err2)
	}
	removeConf := types.ContainerRemoveOptions{
		ContainerID:   c.ID,
		RemoveVolumes: false,
	}
	cli.ContainerRemove(removeConf)

}

// ReadReleaseStatusFromContainer read release.json file from workspace volume container
func ReadReleaseStatusFromContainer(project string) (rs ermtypes.ReleaseStep, err error) {
var releaseStep ermtypes.ReleaseStep

// Get a handle to the KV API
kv := consulCli.KV()

// Lookup the pair
pair, _, err := kv.Get(project, nil)
if err != nil {
    panic(err)
}
fmt.Printf("KV: %v", pair)
buf := bytes.NewBuffer(pair.Value)

decoded, err := base64.StdEncoding.DecodeString( buf.String())
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	fmt.Println(string(decoded))

/*
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var releaseStep ermtypes.ReleaseStep
	var volumeName = project + "-workspace"
	var hostConf = container.HostConfig{

		Binds: []string{volumeName + ":/workspace"},
	}
	var conf = container.Config{
		Image: "busybox",
		Cmd:   []string{"cat", "/workspace/release.json"},
	}

	c, err := cli.ContainerCreate(&conf, &hostConf, nil, "")
	if err != nil {
		return releaseStep, err
	}
	err2 := cli.ContainerStart(c.ID)
	if err2 != nil {
		return releaseStep, err2
	}
	reader, err := cli.ContainerLogs(ctx, types.ContainerLogsOptions{
		ContainerID: c.ID,
		ShowStdout:  true,
	})
	if err != nil {
		log.Print("error ")
		log.Fatal(err)
	}

	if reader != nil {
		defer reader.Close()
		content := readLines(bufio.NewReader(reader))
		log.Println(content)

		var release ermtypes.Release
		dec := json.NewDecoder(reader)
		errJson := dec.Decode(&release)
		if errJson != nil {
			log.Print("release.json to JSON error ")
			log.Fatal(errJson)
		}
		releaseStep = release.ReleaseStep
	}
	/*
	   buf := new(bytes.Buffer)
	   s := buf.String()
	   	_, err = io.Copy(content, reader)
	   	if err != nil && err != io.EOF {
	   		log.Fatal(err)
	   	}
	*/

	return releaseStep, err

}

func readLines(b *bufio.Reader) string {

	s := ""
	for {
		s1, err := b.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic("GetLines: " + err.Error())
		}
		s += s1
	}
	return s
}

// ListenToDockerEvents useful to listen to Docker Events
// Need to update it with engine-api
func ListenToDockerEvents() {
	// Init the Docker client
	var tlsc tls.Config
	cert, err := tls.LoadX509KeyPair(os.Getenv("DOCKER_CERT_PATH")+"/cert.pem", os.Getenv("DOCKER_CERT_PATH")+"/key.pem")
	tlsc.Certificates = append(tlsc.Certificates, cert)
	tlsc.InsecureSkipVerify = true
	docker, err := dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), &tlsc)
	if err != nil {
		log.Fatal(err)
	}
	// Listen to events
	docker.StartMonitorEvents(eventCallback, nil)
}

// Callback used to listen to Docker's events
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	// log.Printf("Received event: %#v\n", *event)

	//TODO: Start a new docker container when a release is OK in a previous container
}

// ListContainers Display running containers
func ListContainers(running bool) {
	options := types.ContainerListOptions{All: running}
	containers, err := cli.ContainerList(options)
	if err != nil {
		panic(err)
	}

    if (len(containers) == 0){
        fmt.Println("No containers available.")
    }
	for _, c := range containers {
		fmt.Println(c.Names)
	}
}
