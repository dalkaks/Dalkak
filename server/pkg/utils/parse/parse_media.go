package parseutil

func ConvertKeyToStaticLink(staticLink, key string) string {
	return staticLink + key
}

func ConvertStaticLinkToKey(staticLink, url string) string {
	return url[len(staticLink):]
}
