package acceptance_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2/vacation"
	"github.com/sclevine/spec"

	. "github.com/paketo-buildpacks/jam/integration/matchers"
)

func testMetadataStaticStack(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		tmpDir string
	)

	it.Before(func() {
		var err error
		tmpDir, err = os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	it("builds static stack", func() {
		var runReleaseDate time.Time

		by("confirming that the run image is correct", func() {
			dir := filepath.Join(tmpDir, "run-index")
			err := os.Mkdir(dir, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			archive, err := os.Open(staticStack.RunArchive)
			Expect(err).NotTo(HaveOccurred())
			defer archive.Close()

			err = vacation.NewArchive(archive).Decompress(dir)
			Expect(err).NotTo(HaveOccurred())

			path, err := layout.FromPath(dir)
			Expect(err).NotTo(HaveOccurred())

			index, err := path.ImageIndex()
			Expect(err).NotTo(HaveOccurred())

			indexManifest, err := index.IndexManifest()
			Expect(err).NotTo(HaveOccurred())

			Expect(indexManifest.Manifests).To(HaveLen(2))
			platforms := []v1.Platform{}
			for _, manifest := range indexManifest.Manifests {
				platforms = append(platforms, v1.Platform{
					Architecture: manifest.Platform.Architecture,
					OS:           manifest.Platform.OS,
				})
			}
			Expect(platforms).To(ContainElement(v1.Platform{
				OS:           "linux",
				Architecture: "amd64",
			}))
			Expect(platforms).To(ContainElement(v1.Platform{
				OS:           "linux",
				Architecture: "arm64",
			}))

			image, err := index.Image(indexManifest.Manifests[0].Digest)
			Expect(err).NotTo(HaveOccurred())

			file, err := image.ConfigFile()
			Expect(err).NotTo(HaveOccurred())

			Expect(file.Config.Labels).To(SatisfyAll(
				HaveKeyWithValue("io.buildpacks.stack.id", "io.buildpacks.stacks.noble.static"),
				HaveKeyWithValue("io.buildpacks.stack.description", "noble for statically-linked applications containing tzdata and ca-certificates"),
				HaveKeyWithValue("io.buildpacks.stack.distro.name", "ubuntu"),
				HaveKeyWithValue("io.buildpacks.stack.distro.version", "24.04"),
				HaveKeyWithValue("io.buildpacks.stack.homepage", "https://github.com/paketo-buildpacks/noble-static-stack"),
				HaveKeyWithValue("io.buildpacks.stack.maintainer", "Paketo Buildpacks"),
				HaveKeyWithValue("io.buildpacks.stack.metadata", MatchJSON("{}")),
			))

			runReleaseDate, err = time.Parse(time.RFC3339, file.Config.Labels["io.buildpacks.stack.released"])
			Expect(err).NotTo(HaveOccurred())
			Expect(runReleaseDate).NotTo(BeZero())

			Expect(file.Config.User).To(Equal("1002:1001"))

			Expect(image).To(SatisfyAll(
				HaveFileWithContent("/etc/group", ContainSubstring("cnb:x:1001:")),
				HaveFileWithContent("/etc/passwd", ContainSubstring("cnb:x:1002:1001::/home/cnb:/sbin/nologin")),
				HaveDirectory("/home/cnb"),
			))

			Expect(image).To(SatisfyAll(
				HaveFile("/usr/share/doc/ca-certificates/copyright"),
				HaveFile("/etc/ssl/certs/ca-certificates.crt"),
				HaveDirectory("/tmp"),
				HaveFile("/etc/nsswitch.conf"),
			))

			Expect(image).To(HaveFileWithContent("/var/lib/dpkg/status.d/ca-certificates", SatisfyAll(
				ContainSubstring("Package: ca-certificates"),
				MatchRegexp("Version: [0-9]+"),
				ContainSubstring("Architecture: all"),
			)))

			Expect(image).To(HaveFileWithContent("/var/lib/dpkg/status.d/tzdata", SatisfyAll(
				ContainSubstring("Package: tzdata"),
				MatchRegexp("Version: [a-z0-9\\.\\-]+ubuntu[0-9\\.]+"),
				ContainSubstring("Architecture: all"),
			)))

			Expect(image).NotTo(HaveFile("/usr/share/ca-certificates"))

			Expect(image).To(HaveFileWithContent("/etc/os-release", SatisfyAll(
				ContainSubstring(`PRETTY_NAME="Paketo Buildpacks Static Noble"`),
				ContainSubstring(`HOME_URL="https://github.com/paketo-buildpacks/noble-static-stack"`),
				ContainSubstring(`SUPPORT_URL="https://github.com/paketo-buildpacks/noble-static-stack/blob/main/README.md"`),
				ContainSubstring(`BUG_REPORT_URL="https://github.com/paketo-buildpacks/noble-static-stack/issues/new"`),
			)))
		})
	})
}
