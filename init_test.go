package acceptance_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/onsi/gomega/format"
	"github.com/paketo-buildpacks/occam"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

type Builder struct {
	LocalInfo struct {
		Lifecycle struct {
			Version string `json:"version"`
		} `json:"lifecycle"`
	} `json:"local_info"`
}

var tinyStack struct {
	BuildArchive string
	RunArchive   string
	BuildImageID string
	RunImageID   string
}

var baseStack struct {
	BuildArchive string
	RunArchive   string
	BuildImageID string
	RunImageID   string
}

var RegistryUrl string

var lifecycleVersion string

func by(_ string, f func()) { f() }

func TestAcceptance(t *testing.T) {
	docker := occam.NewDocker()

	format.MaxLength = 0
	SetDefaultEventuallyTimeout(30 * time.Second)

	Expect := NewWithT(t).Expect

	RegistryUrl = os.Getenv("REGISTRY_URL")
	Expect(RegistryUrl).NotTo(Equal(""))

	root, err := filepath.Abs(".")
	Expect(err).ToNot(HaveOccurred())

	tinyStack.BuildArchive = filepath.Join(root, "builds", "noble-tiny-stack", "build.oci")
	tinyStack.BuildImageID = fmt.Sprintf("%s/noble-tiny-stack-build-%s", RegistryUrl, uuid.NewString())

	tinyStack.RunArchive = filepath.Join(root, "builds", "noble-tiny-stack", "run.oci")
	tinyStack.RunImageID = fmt.Sprintf("%s/noble-tiny-stack-run-%s", RegistryUrl, uuid.NewString())

	baseStack.BuildArchive = filepath.Join(root, "builds", "noble-base-stack", "build.oci")
	baseStack.BuildImageID = fmt.Sprintf("%s/noble-base-stack-build-%s", RegistryUrl, uuid.NewString())

	baseStack.RunArchive = filepath.Join(root, "builds", "noble-base-stack", "run.oci")
	baseStack.RunImageID = fmt.Sprintf("%s/noble-base-stack-run-%s", RegistryUrl, uuid.NewString())

	suite := spec.New("Acceptance", spec.Report(report.Terminal{}), spec.Parallel())
	suite("MetadataTinyStack", testMetadataTinyStack)
	suite("MetadataBaseStack", testMetadataBaseStack)
	suite("BuildpackIntegrationTinyStack", testBuildpackIntegrationTinyStack)
	suite("BuildpackIntegrationBaseStack", testBuildpackIntegrationBaseStack)
	suite.Run(t)

	Expect(docker.Image.Remove.Execute(fmt.Sprintf("buildpacksio/lifecycle:%s", lifecycleVersion))).To(Succeed())

}

func createBuilder(config string, name string) (string, error) {
	buf := bytes.NewBuffer(nil)

	pack := pexec.NewExecutable("pack")
	err := pack.Execute(pexec.Execution{
		Stdout: buf,
		Stderr: buf,
		Args: []string{
			"builder",
			"create",
			name,
			fmt.Sprintf("--config=%s", config),
		},
	})
	return buf.String(), err
}

func getLifecycleVersion(builderID string) (string, error) {
	buf := bytes.NewBuffer(nil)
	pack := pexec.NewExecutable("pack")
	err := pack.Execute(pexec.Execution{
		Stdout: buf,
		Stderr: buf,
		Args: []string{
			"builder",
			"inspect",
			builderID,
			"-o",
			"json",
		},
	})

	if err != nil {
		return "", err
	}

	var builder Builder
	err = json.Unmarshal(buf.Bytes(), &builder)
	if err != nil {
		return "", err
	}

	lifecycleVersion = builder.LocalInfo.Lifecycle.Version
	return builder.LocalInfo.Lifecycle.Version, nil
}
