package pack_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/fatih/color"

	"github.com/buildpack/pack/config"

	"github.com/buildpack/pack/logging"

	"github.com/buildpack/lifecycle"
	"github.com/golang/mock/gomock"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpack/pack"
	"github.com/buildpack/pack/mocks"
	h "github.com/buildpack/pack/testhelpers"
)

func TestRebase(t *testing.T) {
	color.NoColor = true
	spec.Run(t, "rebase_factory", testRebase, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testRebase(t *testing.T, when spec.G, it spec.S) {
	when("#Rebase", func() {
		var (
			mockController *gomock.Controller
			MockImageFetcher    *mocks.MockImageFetcher
			factory        pack.RebaseFactory
			outBuf         bytes.Buffer
			errBuff        bytes.Buffer
		)
		it.Before(func() {
			mockController = gomock.NewController(t)
			MockImageFetcher = mocks.NewMockImageFetcher(mockController)

			factory = pack.RebaseFactory{
				Logger:  logging.NewLogger(&outBuf, &errBuff, false, false),
				Fetcher: MockImageFetcher,
				Config: &config.Config{},
			}
		})

		it.After(func() {
			mockController.Finish()
		})

		when("#RebaseConfigFromFlags", func() {
			when("run image is provided by the user", func() {
				when("the image has a label with a run image specified", func() {
					it("uses the run image provided by the user", func() {
						mockBaseImage := mocks.NewMockImage(mockController)
						mockBaseImage.EXPECT().Name().Return("some/base-image").AnyTimes()
						mockImage := mocks.NewMockImage(mockController)
						mockImage.EXPECT().Name().Return("some/name").AnyTimes()
						mockImage.EXPECT().Label("io.buildpacks.lifecycle.metadata").
							Return(`{"stack":{"runImage":{"image":"label/run/image"}}}`, nil).AnyTimes()
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "myorg/myrepo", true, true).Return(mockImage, nil)
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "my/run/image", true, true).Return(mockBaseImage, nil)

						flags := pack.RebaseFlags{
							RunImage: "my/run/image",
							RepoName: "myorg/myrepo",
						}

						_, err := factory.RebaseConfigFromFlags(context.TODO(), flags)
						h.AssertNil(t, err)
					})
				})
				when("the image does not have a label with a run image specified", func() {
					it("uses the run image provided by the user", func() {
						mockBaseImage := mocks.NewMockImage(mockController)
						mockBaseImage.EXPECT().Name().Return("some/base-image").AnyTimes()
						mockImage := mocks.NewMockImage(mockController)
						mockImage.EXPECT().Name().Return("some/name").AnyTimes()
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "myorg/myrepo", true, true).Return(mockImage, nil)
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "my/run/image", true, true).Return(mockBaseImage, nil)

						flags := pack.RebaseFlags{
							RunImage: "my/run/image",
							RepoName: "myorg/myrepo",
						}

						_, err := factory.RebaseConfigFromFlags(context.TODO(), flags)
						h.AssertNil(t, err)
					})
				})
			})
			when("run image is NOT provided by the user", func() {
				when("the image has a label with a run image specified", func() {
					it("uses the run image provided in the App image label", func() {
						mockBaseImage := mocks.NewMockImage(mockController)
						mockBaseImage.EXPECT().Name().Return("some/base-image").AnyTimes()
						mockImage := mocks.NewMockImage(mockController)
						mockImage.EXPECT().Name().Return("some/name").AnyTimes()
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "myorg/myrepo", true, true).Return(mockImage, nil)
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "some/other/runimage", true, true).Return(mockBaseImage, nil)
						mockImage.EXPECT().Label("io.buildpacks.lifecycle.metadata").
							Return(`{"stack":{"runImage":{"image":"some/other/runimage"}}}`, nil).AnyTimes()

						flags := pack.RebaseFlags{
							RepoName: "myorg/myrepo",
						}

						rc, err := factory.RebaseConfigFromFlags(context.TODO(), flags)
						h.AssertNil(t, err)
						h.AssertSameInstance(t, rc.NewBaseImage, mockBaseImage)
					})
				})

				when("the image has a label with a run image mirrors specified", func() {
					it("chooses the best mirror", func() {
						mockBaseImage := mocks.NewMockImage(mockController)
						mockBaseImage.EXPECT().Name().Return("some/base-image").AnyTimes()
						mockImage := mocks.NewMockImage(mockController)
						mockImage.EXPECT().Name().Return("some/name").AnyTimes()
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "example.com/myorg/myrepo", true, true).Return(mockImage, nil)
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "example.com/local/run/image", true, true).Return(mockBaseImage, nil)
						mockImage.EXPECT().Label("io.buildpacks.lifecycle.metadata").
							Return(`{"stack":{"runImage":{"image":"some/other/runimage", "mirrors":["example.com/run/image"]}}}`, nil).AnyTimes()
						factory.Config = &config.Config{
							RunImages: []config.RunImage{
								{Image: "some/other/runimage", Mirrors: []string{"example.com/local/run/image"}},
							},
						}
						flags := pack.RebaseFlags{
							RepoName: "example.com/myorg/myrepo",
						}

						rc, err := factory.RebaseConfigFromFlags(context.TODO(), flags)
						h.AssertNil(t, err)
						h.AssertSameInstance(t, rc.NewBaseImage, mockBaseImage)
					})
				})

				when("there are locally configured mirrors", func() {
					it("chooses the best mirror from local and label mirrors", func() {
						mockBaseImage := mocks.NewMockImage(mockController)
						mockBaseImage.EXPECT().Name().Return("some/base-image").AnyTimes()
						mockImage := mocks.NewMockImage(mockController)
						mockImage.EXPECT().Name().Return("some/name").AnyTimes()
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "example.com/myorg/myrepo", true, true).Return(mockImage, nil)
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "example.com/run/image", true, true).Return(mockBaseImage, nil)
						mockImage.EXPECT().Label("io.buildpacks.lifecycle.metadata").
							Return(`{"stack":{"runImage":{"image":"some/other/runimage", "mirrors":["example.com/run/image"]}}}`, nil).AnyTimes()

						flags := pack.RebaseFlags{
							RepoName: "example.com/myorg/myrepo",
						}

						rc, err := factory.RebaseConfigFromFlags(context.TODO(), flags)
						h.AssertNil(t, err)
						h.AssertSameInstance(t, rc.NewBaseImage, mockBaseImage)
					})
				})

				when("the image does not have a label with a run image specified", func() {
					it("returns an error", func() {
						mockImage := mocks.NewMockImage(mockController)
						mockImage.EXPECT().Name().Return("some/name").AnyTimes()
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "myorg/myrepo", true, true).Return(mockImage, nil)
						mockImage.EXPECT().Label("io.buildpacks.lifecycle.metadata").Return(`{"stack":{}}`, nil)

						flags := pack.RebaseFlags{
							RepoName: "myorg/myrepo",
						}

						_, err := factory.RebaseConfigFromFlags(context.TODO(), flags)
						h.AssertError(t, err, "run image must be specified")
					})
				})
			})

			when("publish is false", func() {
				when("no-pull is false", func() {
					it("XXXX", func() {
						mockBaseImage := mocks.NewMockImage(mockController)
						mockBaseImage.EXPECT().Name().Return("some/base-image").AnyTimes()
						mockImage := mocks.NewMockImage(mockController)
						mockImage.EXPECT().Name().Return("some/name").AnyTimes()
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "myorg/myrepo", true, true).Return(mockImage, nil)
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "default/run", true, true).Return(mockBaseImage, nil)

						cfg, err := factory.RebaseConfigFromFlags(context.TODO(), pack.RebaseFlags{
							RepoName: "myorg/myrepo",
							RunImage: "default/run",
							Publish:  false,
							NoPull:   false,
						})
						h.AssertNil(t, err)

						h.AssertSameInstance(t, cfg.Image, mockImage)
						h.AssertSameInstance(t, cfg.NewBaseImage, mockBaseImage)
					})
				})

				when("no-pull is true", func() {
					it("XXXX", func() {
						mockBaseImage := mocks.NewMockImage(mockController)
						mockBaseImage.EXPECT().Name().Return("some/base-image").AnyTimes()
						mockImage := mocks.NewMockImage(mockController)
						mockImage.EXPECT().Name().Return("some/name").AnyTimes()
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "myorg/myrepo", true, false).Return(mockImage, nil)
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "default/run", true, false).Return(mockBaseImage, nil)

						cfg, err := factory.RebaseConfigFromFlags(context.TODO(), pack.RebaseFlags{
							RepoName: "myorg/myrepo",
							RunImage: "default/run",
							Publish:  false,
							NoPull:   true,
						})
						h.AssertNil(t, err)

						h.AssertSameInstance(t, cfg.Image, mockImage)
						h.AssertSameInstance(t, cfg.NewBaseImage, mockBaseImage)
					})
				})
			})

			when("publish is true", func() {
				when("no-pull is anything", func() {
					it("XXXX", func() {
						mockBaseImage := mocks.NewMockImage(mockController)
						mockBaseImage.EXPECT().Name().Return("some/base-image").AnyTimes()
						mockImage := mocks.NewMockImage(mockController)
						mockImage.EXPECT().Name().Return("some/name").AnyTimes()
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "myorg/myrepo", false, true).Return(mockImage, nil)
						MockImageFetcher.EXPECT().Fetch(gomock.Any(), "default/run", false, true).Return(mockBaseImage, nil)

						cfg, err := factory.RebaseConfigFromFlags(context.TODO(), pack.RebaseFlags{
							RepoName: "myorg/myrepo",
							RunImage: "default/run",
							Publish:  true,
							NoPull:   false,
						})
						h.AssertNil(t, err)

						h.AssertSameInstance(t, cfg.Image, mockImage)
						h.AssertSameInstance(t, cfg.NewBaseImage, mockBaseImage)
					})
				})
			})
		})

		when("#Rebase", func() {
			it("swaps the old base for the new base AND stores new sha for new runimage", func() {
				mockBaseImage := mocks.NewMockImage(mockController)
				mockBaseImage.EXPECT().Name().Return("some/base-image").AnyTimes()
				mockBaseImage.EXPECT().TopLayer().Return("some-top-layer", nil)
				mockBaseImage.EXPECT().Digest().Return("some-sha", nil)
				mockImage := mocks.NewMockImage(mockController)
				mockImage.EXPECT().Name().Return("some/name").AnyTimes()
				mockImage.EXPECT().Label("io.buildpacks.lifecycle.metadata").
					Return(`{"runimage":{"topLayer":"old-top-layer"}, "app":{"sha":"data"}}`, nil)
				mockImage.EXPECT().Rebase("old-top-layer", mockBaseImage)
				setLabel := mockImage.EXPECT().SetLabel("io.buildpacks.lifecycle.metadata", gomock.Any()).
					Do(func(_, label string) {
						var metadata lifecycle.AppImageMetadata
						h.AssertNil(t, json.Unmarshal([]byte(label), &metadata))
						h.AssertEq(t, metadata.RunImage.TopLayer, "some-top-layer")
						h.AssertEq(t, metadata.RunImage.SHA, "some-sha")
						h.AssertEq(t, metadata.App.SHA, "data")
					})
				mockImage.EXPECT().Save().After(setLabel).Return("some-digest", nil)

				rebaseConfig := pack.RebaseConfig{
					Image:        mockImage,
					NewBaseImage: mockBaseImage,
				}
				err := factory.Rebase(rebaseConfig)
				h.AssertNil(t, err)
			})
		})
	})
}
