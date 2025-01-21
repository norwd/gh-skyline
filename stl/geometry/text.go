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
	baseWidthVoxelResolution = 2000 // Number of voxels across the skyline face
	voxelDepth      = 1.0 // Distance to come out of face

	logoScale      = 0.4 // Percent
	logoTopOffset  = 0.15 // Percent
	logoLeftOffset = 0.03 // Percent
	
	usernameFontSize     = 120.0
	usernameLeftOffset   = 0.1 // Percent
	
	yearFontSize         = 100.0
	yearLeftOffset       = 0.85 // Percent
)

// Create3DText generates 3D text geometry for the username and year.
func Create3DText(username string, year string, baseWidth float64, baseHeight float64) ([]types.Triangle, error) {
	if username == "" {
		username = "anonymous"
	}

	usernameTriangles, err := renderText(
		username,
		usernameLeftOffset,
		usernameFontSize,
		baseWidth,
		baseHeight,
	)
	if err != nil {
		return nil, err
	}

	yearTriangles, err := renderText(
		year,
		yearLeftOffset,
		yearFontSize,
		baseWidth,
		baseHeight,
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
//   fontSize (float64): How large to make the text. Note: It scales with the baseWidthVoxelResolution.
//
// Returns:
//   ([]types.Triangle, error): A slice of triangles representing text.
func renderText(text string, leftOffsetPercent float64, fontSize float64, baseWidth float64, baseHeight float64) ([]types.Triangle, error) {
	// Create a rendering context for the face of the skyline
	faceWidthRes := baseWidthVoxelResolution
	faceHeightRes := int(float64(faceWidthRes) * baseHeight/baseWidth)
	
	// Create image representing the skyline face
	dc := gg.NewContext(faceWidthRes, faceHeightRes)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)

	// Load font into context
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

	// Draw text on image at desired location
	var triangles []types.Triangle

	// Draw pixelated text in image at desired location
	dc.DrawStringAnchored(
		text,
		float64(faceWidthRes)*leftOffsetPercent, // Offset from right
		float64(faceHeightRes)*0.5, // Offset from top
		0.0, // Left aligned
		0.5, // Vertically aligned
	)

	// Convert context image pixels into voxels
	for x := 0; x < faceWidthRes; x++ {
		for y := 0; y < faceHeightRes; y++ {
			if isPixelActive(dc, x, y) {
				voxel, err := createVoxelOnFace(
					float64(x),
					float64(y),
					voxelDepth,
					baseWidth,
					baseHeight,
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
func createVoxelOnFace(x float64, y float64, height float64, baseWidth float64, baseHeight float64) ([]types.Triangle, error) {
	// Mapping resolution
	xResolution := float64(baseWidthVoxelResolution)
	yResolution := xResolution * baseHeight / baseWidth

	// Pixel size
	voxelSize := 1.0;

	// Scale coordinate to face resolution
	x = (x / xResolution) * baseWidth
	y = (y / yResolution) * baseHeight
	voxelSizeX := (voxelSize / xResolution) * baseWidth;
	voxelSizeY := (voxelSize / yResolution) * baseHeight;

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
func GenerateImageGeometry(baseWidth float64, baseHeight float64) ([]types.Triangle, error) {
	// Get temporary image file
	imgPath, cleanup, err := getEmbeddedImage()
	if err != nil {
		return nil, err
	}

	defer cleanup()

	return renderImage(
		imgPath,
		logoScale,
		voxelDepth,
		logoLeftOffset,
		logoTopOffset,
		baseWidth,
		baseHeight,
	)
}

// renderImage generates 3D geometry for the given image configuration.
// func renderImage(config imageRenderConfig) ([]types.Triangle, error) {
func renderImage(filePath string, scale float64, height float64, leftOffsetPercent float64, topOffsetPercent float64, baseWidth float64, baseHeight float64) ([]types.Triangle, error) {

	// Get voxel resolution of base face
	faceWidthRes := baseWidthVoxelResolution
	faceHeightRes := int(float64(faceWidthRes) * baseHeight/baseWidth)
	
	// Load image from file
	reader, err := os.Open(filePath)
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

	// Get image size
	bounds := img.Bounds()
	logoWidth := bounds.Max.X
	logoHeight := bounds.Max.Y

	// Transfer image pixels onto face of skyline as voxels
	var triangles []types.Triangle
	for x := 0; x < logoWidth; x++ {
		for y := logoHeight - 1; y >= 0; y-- {
			// Get pixel color and alpha
			r, _, _, a := img.At(x, y).RGBA()

			// If pixel is active (white) and not fully transparent, create a voxel
			if a > 32768 && r > 32768 {

				voxel, err := createVoxelOnFace(
					(leftOffsetPercent * float64(faceWidthRes)) + float64(x)*logoScale,
					(topOffsetPercent * float64(faceHeightRes)) + float64(y)*logoScale,
					height,
					baseWidth,
					baseHeight,
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
