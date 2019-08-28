package utils

func PathParser(dir, defaultPath string) string {
	if "" == dir {
		if "" != defaultPath {
			return PathParser(defaultPath, "")
		} else {
			return defaultPath
		}
	}
	if dir[:1] != "/" {
		dir = "/" + dir
	}
	if dir[len(dir)-1:] == "/" {
		dir = dir[:len(dir)-1]
	}
	return dir
}
