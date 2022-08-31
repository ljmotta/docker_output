package pkg

import (
	"bufio"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "example",
	}
	cmd.AddCommand(NewBuildCommand())
	return cmd
}

func NewBuildCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "build",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runBuild(cmd)
	}
	return cmd
}

func runBuild(cmd *cobra.Command) (err error) {
	var dockerClient client.CommonAPIClient
	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	tar, err := archive.TarWithOptions("./", &archive.TarOptions{})
	if err != nil {
		return err
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile.test",
		Tags:       []string{"test"},
		Version:    types.BuilderBuildKit,
		PullParent: true,
		Target:     "output",
		Outputs: []types.ImageBuildOutput{{
			Type: "local",
			Attrs: map[string]string{
				"dest": "output",
			},
		}},
	}

	res, err := dockerClient.ImageBuild(cmd.Context(), tar, opts)
	if err != nil {
		return fmt.Errorf("cannot build the app image: %w", err)
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return
}
