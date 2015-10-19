package main

import (
	"testing"

	"github.com/bamiaux/rez"
	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
	nfnt_resize "github.com/nfnt/resize"
	x_draw "golang.org/x/image/draw"
)

const runParallel = false

func BenchmarkNfntNearestNeighbor(b *testing.B) {
	benchmarkNfnt(b, nfnt_resize.NearestNeighbor)
}

func BenchmarkNfntBilinear(b *testing.B) {
	benchmarkNfnt(b, nfnt_resize.Bilinear)
}

func BenchmarkNfntBicubic(b *testing.B) {
	benchmarkNfnt(b, nfnt_resize.Bicubic)
}

func BenchmarkNfntMitchellNetravali(b *testing.B) {
	benchmarkNfnt(b, nfnt_resize.MitchellNetravali)
}

func BenchmarkNfntLanczos2(b *testing.B) {
	benchmarkNfnt(b, nfnt_resize.Lanczos2)
}

func BenchmarkNfntLanczos3(b *testing.B) {
	benchmarkNfnt(b, nfnt_resize.Lanczos3)
}

func benchmarkNfnt(b *testing.B, interp nfnt_resize.InterpolationFunction) {
	benchmark(b, newResizeFuncNfnt(interp))
}

func BenchmarkXDrawNearestNeighbor(b *testing.B) {
	benchmarkXDraw(b, x_draw.NearestNeighbor)
}

func BenchmarkXDrawApproxBiLinear(b *testing.B) {
	benchmarkXDraw(b, x_draw.ApproxBiLinear)
}

func BenchmarkXDrawBiLinear(b *testing.B) {
	benchmarkXDraw(b, x_draw.BiLinear)
}

func BenchmarkXDrawCatmullRom(b *testing.B) {
	benchmarkXDraw(b, x_draw.CatmullRom)
}

func benchmarkXDraw(b *testing.B, interp x_draw.Interpolator) {
	benchmark(b, newResizeFuncXDraw(interp))
}

func BenchmarkImagingNearestNeighbor(b *testing.B) {
	benchmarkImaging(b, imaging.NearestNeighbor)
}

func BenchmarkImagingBox(b *testing.B) {
	benchmarkImaging(b, imaging.Box)
}

func BenchmarkImagingLinear(b *testing.B) {
	benchmarkImaging(b, imaging.Linear)
}

func BenchmarkImagingMitchellNetravali(b *testing.B) {
	benchmarkImaging(b, imaging.MitchellNetravali)
}

func BenchmarkImagingCatmullRom(b *testing.B) {
	benchmarkImaging(b, imaging.CatmullRom)
}

func BenchmarkImagingGaussian(b *testing.B) {
	benchmarkImaging(b, imaging.Gaussian)
}

func BenchmarkImagingLanczos(b *testing.B) {
	benchmarkImaging(b, imaging.Lanczos)
}

func benchmarkImaging(b *testing.B, filter imaging.ResampleFilter) {
	benchmark(b, newResizeFuncImaging(filter))
}

func BenchmarkGiftNearestNeighbor(b *testing.B) {
	benchmarkGift(b, gift.NearestNeighborResampling)
}

func BenchmarkGiftBox(b *testing.B) {
	benchmarkGift(b, gift.BoxResampling)
}

func BenchmarkGiftLinear(b *testing.B) {
	benchmarkGift(b, gift.LinearResampling)
}

func BenchmarkGiftCubic(b *testing.B) {
	benchmarkGift(b, gift.CubicResampling)
}

func BenchmarkGiftLanczos(b *testing.B) {
	benchmarkGift(b, gift.LanczosResampling)
}

func benchmarkGift(b *testing.B, resampling gift.Resampling) {
	benchmark(b, newResizeFuncGift(resampling))
}

func BenchmarkRezBilinear(b *testing.B) {
	benchmarkRez(b, rez.NewBilinearFilter())
}

func BenchmarkRezBicubic(b *testing.B) {
	benchmarkRez(b, rez.NewBicubicFilter())
}

func BenchmarkRezLanczos2(b *testing.B) {
	benchmarkRez(b, rez.NewLanczosFilter(2))
}

func BenchmarkRezLanczos3(b *testing.B) {
	benchmarkRez(b, rez.NewLanczosFilter(3))
}

func benchmarkRez(b *testing.B, filter rez.Filter) {
	defer func() {
		err := recover()
		if err != nil {
			b.Skipf("skip: %s", err)
		}
	}()
	benchmark(b, newResizeFuncRez(filter))
}

var benchImage = load()

func benchmark(b *testing.B, res resizeFunc) {
	if runParallel {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				res(benchImage)
			}
		})
	} else {
		for i := 0; i < b.N; i++ {
			res(benchImage)
		}
	}
}
