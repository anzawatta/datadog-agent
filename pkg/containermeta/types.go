// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package containermeta

type Kind string
type ContainerRuntime string
type ECSLaunchType string
type EventType int

const (
	KindContainer     Kind = "container"
	KindKubernetesPod Kind = "kubernetes_pod"
	KindECSTask       Kind = "ecs_task"

	ContainerRuntimeDocker ContainerRuntime = "docker"

	ECSLaunchTypeEC2      ECSLaunchType = "ec2"
	ECSLaunchTypeFargate  ECSLaunchType = "fargate"
	ECSLaunchTypeExternal ECSLaunchType = "external"

	EventTypeSet EventType = iota
	EventTypeUnset
)

type Entity interface {
	GetMeta() EntityMeta
}

type EntityMeta struct {
	Kind        Kind
	ID          string
	Name        string
	Namespace   string
	Annotations map[string]string
	Labels      map[string]string
}

type ContainerImage struct {
	ID        string
	Name      string
	ShortName string
	Tag       string
}

type Container struct {
	EntityMeta
	Image   ContainerImage
	EnvVars map[string]string
	Runtime ContainerRuntime
}

func (c Container) GetMeta() EntityMeta {
	return c.EntityMeta
}

type KubernetesPod struct {
	EntityMeta
	Owners     []interface{}
	Containers []string
	Phase      string
}

func (p KubernetesPod) GetMeta() EntityMeta {
	return p.EntityMeta
}

type ECSTask struct {
	EntityMeta
	Containers []Container
	LaunchType ECSLaunchType
}

func (t ECSTask) GetMeta() EntityMeta {
	return t.EntityMeta
}

type Event struct {
	Type   EventType
	Entity Entity
}
