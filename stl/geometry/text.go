package geometry

import (
	"fmt"
	"image/png"
	"os"

	"github.com/fogleman/gg"
	"github.com/github/gh-skyline/errors"
	"github.com/github/gh-skyline/types"
)

// Common configuration for rendered elements
type renderConfig struct {
	startX     float64
	startY     float64
	startZ     float64
	voxelScale float64
	depth      float64
}

// ImageConfig holds parameters for image rendering
type imageRenderConfig struct {
	renderConfig
	imagePath string
	height    float64
}

const (
	skylineFaceWidth  = 142.5 // millimeters?
	skylineFaceHeight = 10.0 // millimeters?
	skylineResolutionWidth = 2000 // voxels

	logoPosition  = 0.025
	logoHeight = 9.0
	logoScale  = 0.8
	logoLeftMargin    = 10.0
	
	voxelDepth      = 1.0
	
	usernameFontSize     = 120.0
	usernameLeftOffset   = 0.1
	
	yearFontSize         = 100.0
	yearLeftOffset       = 0.85
)

// Create3DText generates 3D text geometry for the username and year.
func Create3DText(username string, year string, innerWidth, baseHeight float64) ([]types.Triangle, error) {
	if username == "" {
		username = "anonymous"
	}

	usernameTriangles, err := renderText(
		username,
		usernameLeftOffset,
		usernameFontSize,
	)
	if err != nil {
		return nil, err
	}

	yearTriangles, err := renderText(
		year,
		yearLeftOffset,
		yearFontSize,
	)
	if err != nil {
		return nil, err
	}

	return append(usernameTriangles, yearTriangles...), nil
}

// renderText places text on the face of a skyline, offset from the left and vertically-aligned.
// The function takes the text to be displayed, offset from left, and font size.
// It returns an array of types.Triangle.
//
// Parameters:
//   text (string): The text to be displayed on the skyline's front face.
//   leftOffsetPercent (float64): The percentage distance from the left to start displaying the text.
//   fontSize (float64): How large to make the text. Note: It scales with the skylineResolutionWidth.
//
// Returns:
//   ([]types.Triangle, error): A slice of triangles representing text.
func renderText(text string, leftOffsetPercent float64, fontSize float64) ([]types.Triangle, error) {
	// Create a rendering context for the face of the skyline
	faceWidthRes := skylineResolutionWidth
	faceHeightRes := int(float64(faceWidthRes) * skylineFaceHeight/skylineFaceWidth)
	dc := gg.NewContext(faceWidthRes, faceHeightRes)

	// Get temporary font file
	fontPath, cleanup, err := writeTempFont(PrimaryFont)
	if err != nil {
		// Try fallback font
		fontPath, cleanup, err = writeTempFont(FallbackFont)
		if err != nil {
			return nil, errors.New(errors.IOError, "failed to load any fonts", err)
		}
	}

	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return nil, errors.New(errors.IOError, "failed to load font", err)
	}

	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)

	var triangles []types.Triangle

	// Draw pixelated text in image at desired location
	dc.DrawStringAnchored(
		text,
		float64(faceWidthRes)*leftOffsetPercent, // Offset from left
		float64(faceHeightRes)*0.5, // Offset from top
		0.0, // Left aligned
		0.5, // Vertically aligned
	)

	// Transfer image pixels onto face of skyline as voxels
	for x := 0; x < faceWidthRes; x++ {
		for y := 0; y < faceHeightRes; y++ {
			if isPixelActive(dc, x, y) {
				voxel, err := createVoxelOnFace(
					float64(x),
					float64(y),
					voxelDepth,
				)
				if err != nil {
					return nil, errors.New(errors.STLError, "failed to create cube", err)
				}

				triangles = append(triangles, voxel...)
			}
		}
	}

	defer cleanup()

	return triangles, nil
}

// createVoxelOnFace creates a voxel on the face of a skyline by generating a cube at the specified coordinates.
// The function takes in the x, y coordinates and height.
// It returns a slice of types.Triangle representing the cube and an error if the cube creation fails.
//
// Parameters:
//   x (float64): The x-coordinate on the skyline face (left to right).
//   y (float64): The y-coordinate on the skyline face (top to bottom).
//   height (float64): Distance coming out of the face.
//
// Returns:
//   ([]types.Triangle, error): A slice of triangles representing the cube and an error if any.
func createVoxelOnFace(x float64, y float64, height float64) ([]types.Triangle, error) {
	// Mapping resolution
	xResolution := float64(skylineResolutionWidth)
	yResolution := xResolution * skylineFaceHeight / skylineFaceWidth

	// Pixel size
	voxelSize := 1.0;

	// Scale coordinate to face resolution
	x = (x / xResolution) * skylineFaceWidth
	y = (y / yResolution) * skylineFaceHeight
	voxelSizeX := (voxelSize / xResolution) * skylineFaceWidth;
	voxelSizeY := (voxelSize / yResolution) * skylineFaceHeight;

	cube, err := CreateCube(
		// Location (from top left corner of skyline face)
		x, // x - Left to right
		-height, // y - Negative comes out of face. Positive goes into face.
		-voxelSizeY - y, // z - Bottom to top
		
		// Size
		voxelSizeX, // x length - left to right from specified point
		height, // thickness - distance coming out of face
		voxelSizeY, // y length - bottom to top from specified point
	)

	return cube, err
}

// GenerateImageGeometry creates 3D geometry from the embedded logo image.
func GenerateImageGeometry(innerWidth, baseHeight float64) ([]types.Triangle, error) {
	// Get temporary image file
	imgPath, cleanup, err := getEmbeddedImage()
	if err != nil {
		return nil, err
	}

	config := imageRenderConfig{
		renderConfig: renderConfig{
			startX:     innerWidth * logoPosition,
			startY:     -voxelDepth / 2.0,
			startZ:     -0.85 * baseHeight,
			voxelScale: logoScale,
			depth:      voxelDepth,
		},
		imagePath: imgPath,
		height:    logoHeight,
	}

	defer cleanup()

	return renderImage(config)
}

// renderImage generates 3D geometry for the given image configuration.
func renderImage(config imageRenderConfig) ([]types.Triangle, error) {
	reader, err := os.Open(config.imagePath)
	if err != nil {
		return nil, errors.New(errors.IOError, "failed to open image", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			closeErr := errors.New(errors.IOError, "failed to close reader", err)
			// Log the error or handle it appropriately
			fmt.Println(closeErr)
		}
	}()

	img, err := png.Decode(reader)
	if err != nil {
		return nil, errors.New(errors.IOError, "failed to decode PNG", err)
	}

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	scale := config.height / float64(height)

	var triangles []types.Triangle

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			r, _, _, a := img.At(x, y).RGBA()
			if a > 32768 && r > 32768 {
				xPos := config.startX + float64(x)*config.voxelScale*scale
				zPos := config.startZ + float64(height-1-y)*config.voxelScale*scale

				voxel, err := CreateCube(
					xPos,
					config.startY,
					zPos,
					config.voxelScale*scale,
					config.depth,
					config.voxelScale*scale,
				)

				if err != nil {
					return nil, errors.New(errors.STLError, "failed to create cube", err)
				}

				triangles = append(triangles, voxel...)
			}
		}
	}

	return triangles, nil
}

// isPixelActive checks if a pixel is active (white) in the given context.
func isPixelActive(dc *gg.Context, x, y int) bool {
	r, _, _, _ := dc.Image().At(x, y).RGBA()
	return r > 32768
}
