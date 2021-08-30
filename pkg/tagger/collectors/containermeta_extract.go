// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package collectors

import (
	"context"
	"fmt"
	"strings"

	"github.com/DataDog/datadog-agent/pkg/containermeta"
	"github.com/DataDog/datadog-agent/pkg/tagger/utils"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes/kubelet"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

func (c *ContainerMetaCollector) processEvents(ctx context.Context, evs []containermeta.Event) {
	tagInfos := []*TagInfo{}

	for _, ev := range evs {
		entity := ev.Entity
		meta := entity.GetMeta()

		switch meta.Kind {
		case containermeta.KindContainer:
			// TODO
		case containermeta.KindKubernetesPod:
			tagInfos = append(tagInfos, c.handleKubePod(ev)...)
		case containermeta.KindECSTask:
			// TODO
		default:
			log.Errorf("cannot handle event for entity %q with kind %q", meta.ID, meta.Kind)
		}
	}

	c.out <- tagInfos
}

func (c *ContainerMetaCollector) handleKubePod(ev containermeta.Event) []*TagInfo {
	tagInfos := []*TagInfo{}

	pod := ev.Entity.(containermeta.KubernetesPod)

	tags := utils.NewTagList()

	tags.AddOrchestrator(kubernetes.PodTagName, pod.Name)
	tags.AddLow(kubernetes.NamespaceTagName, pod.Namespace)
	tags.AddLow("pod_phase", strings.ToLower(pod.Phase))

	for name, value := range pod.Labels {
		switch name {
		case kubernetes.EnvTagLabelKey:
			tags.AddStandard(tagKeyEnv, value)
		case kubernetes.VersionTagLabelKey:
			tags.AddStandard(tagKeyVersion, value)
		case kubernetes.ServiceTagLabelKey:
			tags.AddStandard(tagKeyService, value)
		case kubernetes.KubeAppNameLabelKey:
			tags.AddLow(tagKeyKubeAppName, value)
		case kubernetes.KubeAppInstanceLabelKey:
			tags.AddLow(tagKeyKubeAppInstance, value)
		case kubernetes.KubeAppVersionLabelKey:
			tags.AddLow(tagKeyKubeAppVersion, value)
		case kubernetes.KubeAppComponentLabelKey:
			tags.AddLow(tagKeyKubeAppComponent, value)
		case kubernetes.KubeAppPartOfLabelKey:
			tags.AddLow(tagKeyKubeAppPartOf, value)
		case kubernetes.KubeAppManagedByLabelKey:
			tags.AddLow(tagKeyKubeAppManagedBy, value)
		}

		// utils.AddMetadataAsTags(name, value, c.labelsAsTags, c.globLabels, tags)
	}

	// for name, value := range pod.Metadata.Annotations {
	// 	utils.AddMetadataAsTags(name, value, c.annotationsAsTags, c.globAnnotations, tags)
	// }

	// if podTags, found := extractTagsFromMap(podTagsAnnotation, pod.Metadata.Annotations); found {
	// 	for tagName, values := range podTags {
	// 		for _, val := range values {
	// 			tags.AddAuto(tagName, val)
	// 		}
	// 	}
	// }

	// OpenShift pod annotations
	if dcName, found := pod.Annotations["openshift.io/deployment-config.name"]; found {
		tags.AddLow("oshift_deployment_config", dcName)
	}
	if deployName, found := pod.Annotations["openshift.io/deployment.name"]; found {
		tags.AddOrchestrator("oshift_deployment", deployName)
	}

	for _, _ = range pod.Owners {
		// TODO
	}

	if pod.ID != "" {
		low, orch, high, standard := tags.Compute()
		tagInfos = append(tagInfos, &TagInfo{
			Source:               containermetaCollectorName,
			Entity:               kubelet.PodUIDToTaggerEntityName(pod.ID),
			HighCardTags:         high,
			OrchestratorCardTags: orch,
			LowCardTags:          low,
			StandardTags:         standard,
		})
	}

	for _, containerID := range pod.Containers {
		// Ignore containers that haven't been created yet and have no
		// ID.
		if containerID == "" {
			continue
		}

		container, err := c.store.GetContainer(containerID)
		if err != nil {
			// TODO(juliogreff): handle missing containers
		}

		cTags := tags.Copy()
		cTags.AddLow("kube_container_name", container.Name)
		cTags.AddHigh("container_id", container.ID)
		if container.Name != "" && pod.Name != "" {
			cTags.AddHigh("display_container_name", fmt.Sprintf("%s_%s", container.Name, pod.Name))
		}

		// Enrich with standard tags from labels for this container if present
		labelTags := []string{tagKeyEnv, tagKeyVersion, tagKeyService}
		for _, tag := range labelTags {
			label := fmt.Sprintf(podStandardLabelPrefix+"%s.%s", container.Name, tag)

			if value, ok := pod.Labels[label]; ok {
				cTags.AddStandard(tag, value)
			}
		}

		// container-specific tags provided through pod annotation
		// containerTags, found := extractTagsFromMap(
		// 	fmt.Sprintf(podContainerTagsAnnotationFormat, container.Name),
		// 	pod.Metadata.Annotations,
		// )
		// if found {
		// 	for tagName, values := range containerTags {
		// 		for _, val := range values {
		// 			cTags.AddAuto(tagName, val)
		// 		}
		// 	}
		// }

		for name, value := range container.EnvVars {
			if value == "" {
				continue
			}

			switch name {
			case envVarEnv:
				cTags.AddStandard(tagKeyEnv, value)
			case envVarVersion:
				cTags.AddStandard(tagKeyVersion, value)
			case envVarService:
				cTags.AddStandard(tagKeyService, value)
			}
		}

		image := container.Image
		cTags.AddLow("image_name", image.Name)
		cTags.AddLow("short_image", image.ShortName)

		imageTag := image.Tag
		if imageTag == "" {
			// k8s default to latest if tag is omitted
			imageTag = "latest"
		}

		cTags.AddLow("image_tag", imageTag)

		low, orch, high, standard := cTags.Compute()
		tagInfos = append(tagInfos, &TagInfo{
			Source:               containermetaCollectorName,
			Entity:               kubelet.RawContainerIDToTaggerEntityID(pod.ID),
			HighCardTags:         high,
			OrchestratorCardTags: orch,
			LowCardTags:          low,
			StandardTags:         standard,
		})
	}

	return tagInfos
}

func (c *ContainerMetaCollector) buildEntityID(meta containermeta.EntityMeta) string {
	return ""
}
