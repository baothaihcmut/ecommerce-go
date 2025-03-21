package results

type CreateProductResult struct {
	ProductResult
	UploadThumbnailDetail *UploadImageResult
	UploadImageDetails    []UploadImageResult
}

type UploadImageResult struct {
	Url    string
	Method string
	Expiry int
}
