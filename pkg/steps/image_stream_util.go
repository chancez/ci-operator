package steps

import (
	"encoding/json"
	"fmt"

	imageapi "github.com/openshift/api/image/v1"
	imageclientset "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	coreapi "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createImageStream(
	isClient imageclientset.ImageStreamsGetter,
	is *imageapi.ImageStream,
	dry bool,
) error {
	if dry {
		isJSON, err := json.MarshalIndent(is, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal imagestream: %v", err)
		}
		fmt.Printf("%s\n", isJSON)
		return nil
	}

	_, err := isClient.ImageStreams(is.Namespace).Create(is)
	return err
}
func createImageStreamTag(
	istClient imageclientset.ImageStreamTagsGetter,
	ist *imageapi.ImageStreamTag,
	dry bool,
) error {
	if dry {
		istJSON, err := json.MarshalIndent(ist, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal imagestreamtag: %v", err)
		}
		fmt.Printf("%s\n", istJSON)
		return nil
	}
	// Create if not exists, update if it does
	_, err := istClient.ImageStreamTags(ist.Namespace).Create(ist)
	return err
}

func newImageStream(namespace, name string) *imageapi.ImageStream {
	return &imageapi.ImageStream{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func newImageStreamTag(
	fromNamespace, fromName, fromTag,
	toNamespace, toName, toTag string,
) *imageapi.ImageStreamTag {
	return &imageapi.ImageStreamTag{
		ObjectMeta: meta.ObjectMeta{
			Name:      fmt.Sprintf("%s:%s", toName, toTag),
			Namespace: toNamespace,
		},
		Tag: &imageapi.TagReference{
			ReferencePolicy: imageapi.TagReferencePolicy{
				Type: imageapi.LocalTagReferencePolicy,
			},
			From: &coreapi.ObjectReference{
				Kind:      "ImageStreamImage",
				Name:      fmt.Sprintf("%s@%s", fromName, fromTag),
				Namespace: fromNamespace,
			},
		},
	}
}
