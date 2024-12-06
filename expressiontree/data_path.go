package expressiontree

type DataPath struct {
	MainPath    TDataKey
	ContentPath ContentPath
}

type ContentPath struct {
	Path  string
	Parts []string
}

func (c *ContentPath) IsEmpty() bool {
	return len(c.Parts) == 0
}

func (c *ContentPath) IsSimple() bool {
	return len(c.Parts) == 1
}

func CreateEmptyContentPath() ContentPath {
	return ContentPath{Path: "", Parts: []string{}}
}

func CreateSimpleContentPath(contentPath string) ContentPath {
	return ContentPath{Path: contentPath, Parts: []string{contentPath}}
}

func CreateContentPath(contentPath string, contentPathParts []string) ContentPath {
	return ContentPath{Path: contentPath, Parts: contentPathParts}
}

func CreateDataPathWithMainOnly(mainPath TDataKey) DataPath {
	return DataPath{MainPath: mainPath, ContentPath: CreateEmptyContentPath()}
}

func CreateDataPathWithSimpleContent(mainPath TDataKey, contentPath string) DataPath {
	return DataPath{MainPath: mainPath, ContentPath: CreateSimpleContentPath(contentPath)}
}

func CreateDataPath(mainPath TDataKey, contentPath ContentPath) DataPath {
	return DataPath{MainPath: mainPath, ContentPath: contentPath}
}

/*func CreateDataPath(mainPath TDataKey, contentPath string, contentPathParts []string) DataPath {
	return DataPath{MainPath: mainPath, ContentPath: CreateContentPath(contentPath, contentPathParts)}
}*/
