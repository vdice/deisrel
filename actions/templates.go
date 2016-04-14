package actions

type genParamsComponentAttrs struct {
	Org        string
	Tag        string
	PullPolicy string
}

type genParamsComponentList struct {
	Minio           genParamsComponentAttrs
	Builder         genParamsComponentAttrs
	SlugBuilder     genParamsComponentAttrs
	DockerBuilder   genParamsComponentAttrs
	Controller      genParamsComponentAttrs
	SlugRunner      genParamsComponentAttrs
	Database        genParamsComponentAttrs
	Registry        genParamsComponentAttrs
	WorkflowManager genParamsComponentAttrs
	Logger          genParamsComponentAttrs
	Router          genParamsComponentAttrs
	FluentD         genParamsComponentAttrs
}
