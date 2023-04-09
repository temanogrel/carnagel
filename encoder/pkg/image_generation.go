package encoder

import "context"

type ImageGenerationService interface {
	// InfinityThumbs generates a single 80 image collage for recording thumbs to target path
	// Note, clean up should be handled post calling defer os.RemoveAll on the second return argument
	InfinityThumbs(ctx context.Context, source, target string) error

	// InfinitySprites generates a single grid of images for recording sprites to target path
	// Note, clean up should be handled post calling defer os.RemoveAll on the second return argument
	InfinitySprites(ctx context.Context, source, target string) error

	// InfinityCollage generates a 3x3 grid of images to the target path
	InfinityCollage(ctx context.Context, source, target string) error

	// WordpressCollage generates a 5x5 grid of images with a header to be published
	// on the wordpress sites
	WordpressCollage(ctx context.Context, source, target, filename string) error
}
