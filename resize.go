package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"time"

	"image/draw"

	"github.com/bamiaux/rez"
	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
	nfnt_resize "github.com/nfnt/resize"
	x_draw "golang.org/x/image/draw"
)

const (
	width         = 600
	height        = 400
	convertSource = true
)

var newImageFunc = func(r image.Rectangle) draw.Image {
	return image.NewRGBA(r)
}

func main() {
	im := load()
	resize(im, newResizeFuncNfnt(nfnt_resize.NearestNeighbor), "nfnt_nearest_neighbor.png")
	resize(im, newResizeFuncNfnt(nfnt_resize.Bilinear), "nfnt_bilinear.png")
	resize(im, newResizeFuncNfnt(nfnt_resize.Bicubic), "nfnt_bicubic.png")
	resize(im, newResizeFuncNfnt(nfnt_resize.MitchellNetravali), "nfnt_mitchell_netravali.png")
	resize(im, newResizeFuncNfnt(nfnt_resize.Lanczos2), "nfnt_lanczos2.png")
	resize(im, newResizeFuncNfnt(nfnt_resize.Lanczos3), "nfnt_lanczos3.png")
	resize(im, newResizeFuncXDraw(x_draw.NearestNeighbor), "x_draw_nearest_neighbor.png")
	resize(im, newResizeFuncXDraw(x_draw.ApproxBiLinear), "x_draw_approx_bilinear.png")
	resize(im, newResizeFuncXDraw(x_draw.BiLinear), "x_draw_bilinear.png")
	resize(im, newResizeFuncXDraw(x_draw.CatmullRom), "x_draw_catmull_rom.png")
	resize(im, newResizeFuncImaging(imaging.NearestNeighbor), "imaging_nearest_neighbor.png")
	resize(im, newResizeFuncImaging(imaging.Box), "imaging_box.png")
	resize(im, newResizeFuncImaging(imaging.Linear), "imaging_linear.png")
	resize(im, newResizeFuncImaging(imaging.MitchellNetravali), "imaging_mitchell_netravali.png")
	resize(im, newResizeFuncImaging(imaging.CatmullRom), "imaging_catmull_rom.png")
	resize(im, newResizeFuncImaging(imaging.Gaussian), "imaging_gaussian.png")
	resize(im, newResizeFuncImaging(imaging.Lanczos), "imaging_lanczos.png")
	resize(im, newResizeFuncGift(gift.NearestNeighborResampling), "gift_nearest_neighbor.png")
	resize(im, newResizeFuncGift(gift.BoxResampling), "gift_box.png")
	resize(im, newResizeFuncGift(gift.LinearResampling), "gift_linear.png")
	resize(im, newResizeFuncGift(gift.CubicResampling), "gift_cubic.png")
	resize(im, newResizeFuncGift(gift.LanczosResampling), "gift_lanczos.png")
	func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("error with rez: %s\n", err)
			}
		}()
		resize(im, newResizeFuncRez(rez.NewBilinearFilter()), "rez_bilinear.png")
		resize(im, newResizeFuncRez(rez.NewBicubicFilter()), "rez_bicubic.png")
		resize(im, newResizeFuncRez(rez.NewLanczosFilter(2)), "rez_lanczos2.png")
		resize(im, newResizeFuncRez(rez.NewLanczosFilter(3)), "rez_lanczos3.png")
	}()
}

func load() image.Image {
	f, err := os.Open("image.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	im, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	if convertSource {
		newIm := newImageFunc(image.Rect(0, 0, im.Bounds().Dx(), im.Bounds().Dy()))
		draw.Draw(newIm, newIm.Bounds(), im, im.Bounds().Min, draw.Src)
		im = newIm
	}
	return im
}

func resize(im image.Image, res resizeFunc, fileName string) {
	start := time.Now()
	im = res(im)
	fmt.Printf("%s: %s\n", fileName, time.Since(start))
	save(im, fileName)
}

type resizeFunc func(image.Image) image.Image

func newResizeFuncNfnt(interp nfnt_resize.InterpolationFunction) resizeFunc {
	return func(im image.Image) image.Image {
		return nfnt_resize.Resize(width, height, im, interp)
	}
}

func newResizeFuncXDraw(interp x_draw.Interpolator) resizeFunc {
	return func(im image.Image) image.Image {
		newIm := newImageFunc(image.Rect(0, 0, width, height))
		interp.Scale(newIm, newIm.Bounds(), im, im.Bounds(), x_draw.Src, nil)
		return newIm
	}
}

func newResizeFuncImaging(filter imaging.ResampleFilter) resizeFunc {
	return func(im image.Image) image.Image {
		return imaging.Resize(im, width, height, filter)
	}
}

func newResizeFuncGift(resampling gift.Resampling) resizeFunc {
	return func(im image.Image) image.Image {
		g := gift.New(gift.Resize(width, height, resampling))
		newIm := newImageFunc(g.Bounds(im.Bounds()))
		g.Draw(newIm, im)
		return newIm
	}
}

func newResizeFuncRez(filter rez.Filter) resizeFunc {
	return func(im image.Image) image.Image {
		newIm := newImageFunc(image.Rect(0, 0, width, height))
		err := rez.Convert(newIm, im, filter)
		if err != nil {
			panic(err)
		}
		return newIm
	}
}

func save(im image.Image, fileName string) {
	os.Mkdir("out", 0755)
	f, err := os.Create("out/" + fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, im)
}
