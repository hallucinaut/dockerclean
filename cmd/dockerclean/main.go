package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

type DockerResource struct {
	Type       string
	ID         string
	Name       string
	Size       string
	Created    time.Time
	UsageCount int
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--dry-run" {
		runDryRun()
	} else {
		runCleanup()
	}
}

func runDryRun() {
	fmt.Println(color.CyanString("\n=== DOCKER CLEANUP DRY RUN ===\n"))

	containers := listContainers()
	volumes := listVolumes()
	images := listImages()

	fmt.Printf("Containers to remove: %d\n", len(containers))
	fmt.Printf("Volumes to remove: %d\n", len(volumes))
	fmt.Printf("Images to remove: %d\n\n", len(images))

	if len(containers) > 0 {
		fmt.Println(color.YellowString("Containers:"))
		for _, c := range containers {
			fmt.Printf("  %s - %s\n", c.ID[:12], c.Name)
		}
	}

	if len(volumes) > 0 {
		fmt.Println(color.YellowString("\nVolumes:"))
		for _, v := range volumes {
			fmt.Printf("  %s - %s\n", v.ID[:12], v.Name)
		}
	}

	if len(images) > 0 {
		fmt.Println(color.YellowString("\nImages:"))
		for _, img := range images {
			fmt.Printf("  %s - %s\n", img.ID[:12], img.Name)
		}
	}

	fmt.Println(color.GreenString("\n=== SAFE CLEANUP COMMANDS ==="))
	fmt.Println("docker container prune -f")
	fmt.Println("docker volume prune -f")
	fmt.Println("docker image prune -a -f")
}

func runCleanup() {
	fmt.Println(color.CyanString("\n=== DOCKER CLEANUP ===\n"))

	containers := listContainers()
	volumes := listVolumes()
	images := listImages()

	fmt.Printf("Removing %d containers...\n", len(containers))
	exec.Command("docker", "container", "prune", "-f").Run()

	fmt.Printf("Removing %d volumes...\n", len(volumes))
	exec.Command("docker", "volume", "prune", "-f").Run()

	fmt.Printf("Removing %d images...\n", len(images))
	exec.Command("docker", "image", "prune", "-a", "-f").Run()

	fmt.Println(color.GreenString("\nCleanup complete!"))
}

func listContainers() []DockerResource {
	cmd := exec.Command("docker", "container", "ls", "-a", "--format", "{{.ID}}|{{.Names}}|{{.Status}}|{{.CreatedAt}}")
	output, _ := cmd.CombinedOutput()

	var resources []DockerResource
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 4 {
			resources = append(resources, DockerResource{
				Type:  "container",
				ID:    parts[0],
				Name:  parts[1],
				Size:  parts[2],
			})
		}
	}

	return resources
}

func listVolumes() []DockerResource {
	cmd := exec.Command("docker", "volume", "ls", "--format", "{{.Name}}|{{.Driver}}|{{.Mountpoint}}")
	output, _ := cmd.CombinedOutput()

	var resources []DockerResource
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			resources = append(resources, DockerResource{
				Type: "volume",
				Name: parts[0],
			})
		}
	}

	return resources
}

func listImages() []DockerResource {
	cmd := exec.Command("docker", "image", "ls", "--format", "{{.ID}}|{{.Repository}}|{{.Tag}}|{{.Size}}|{{.CreatedAt}}")
	output, _ := cmd.CombinedOutput()

	var resources []DockerResource
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 5 {
			resources = append(resources, DockerResource{
				Type:   "image",
				ID:     parts[0],
				Name:   fmt.Sprintf("%s:%s", parts[1], parts[2]),
				Size:   parts[3],
			})
		}
	}

	return resources
}