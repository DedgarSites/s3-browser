package datastores

var (
	ImageTypes = map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"bmp":  true,
	}
	VideoTypes = map[string]bool{
		"mp4":  true,
		"mpeg": true,
		"flv":  true,
		"wmv":  true,
	}
)
