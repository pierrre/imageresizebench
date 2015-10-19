package main

import (
	"testing"

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

func BenchmarkXDrawCatmulRom(b *testing.B) {
	benchmarkXDraw(b, x_draw.CatmullRom)
}

func benchmarkXDraw(b *testing.B, interp x_draw.Interpolator) {
	benchmark(b, newResizeFuncXDraw(interp))
}

func benchmark(b *testing.B, res resizeFunc) {
	im := load()
	b.ResetTimer()
	if runParallel {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				res(im)
			}
		})
	} else {

		for i := 0; i < b.N; i++ {
			res(im)
		}
	}
}
