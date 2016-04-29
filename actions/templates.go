package actions

type genParamsComponentAttrs struct {
	Org        string
	Tag        string
	PullPolicy string
}

type genParamsComponentMap map[string]genParamsComponentAttrs

func createParamsComponentMap() genParamsComponentMap {
	return genParamsComponentMap(map[string]genParamsComponentAttrs{})
}
