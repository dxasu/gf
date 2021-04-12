package boot

import "github.com/dxasu/gf/os/gres"

func init() {
	if err := gres.Add("1f8b08000000000002ff0af066661161e060606048cea809604002420c9c0cc9f9796999e9fa104aaf243f37273484958191cbbb3e8e2bba38b5a82cb528964b414141c13125a528b5b818c4b45550b2b2b0b0b050024b04831505e5e79780240a4a9372329395b8b8a2cb3253cb619a5d52d3124b734adc327352418a32f352522bf44a0a7294a0b23999b99925a945c520b315a29554aa957414946a9562b9b8b8b802bcd939b81c64d7b63330308030c23fad68fe6183fb07ec072feffa389066642508cdd6e9a89a05193819208ed787382fa3041a16973cebe33c527372f21541a6859dd19acb06b68c905320a611e1940c2cf152929a5b90935892aa0f0f2bb041a2de50b728a854ebe525e6a6d6821d75255b6e263f0303033f5e477120998bc7598c4c22cc8804831cf42087c1c09246108927f9201b842d7860c1cee070016e1092ab109a4131851cea8228aee06564c0136fc43b221a6e0e564780e2083994518362062303be18c3e70a0e14577c403208c91dac6c1007b3314c666460d06702f100010000ffffc7c852f1d8030000"); err != nil {
		panic(err)
	}
}