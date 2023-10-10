package assets

import (
	_ "embed"

	"nftsiren/pkg/images"
)

var (
	//go:embed images/discord.png
	DiscordLogoData []byte
	//go:embed images/ethereum.png
	EthereumLogoData []byte
	//go:embed images/nftsiren.png
	NftsirenLogoData []byte
	//go:embed images/looksrare.png
	LooksrareLogoData []byte
	//go:embed images/magiceden.png
	MagicedenLogoData []byte
	//go:embed images/opensea.png
	OpenseaLogoData []byte
	//go:embed images/solana.png
	SolanaLogoData []byte
	//go:embed images/twitter.png
	TwitterLogoData []byte
)

var (
	DiscordLogo   = images.MustParse(DiscordLogoData)
	EthereumLogo  = images.MustParse(EthereumLogoData)
	NftsirenLogo  = images.MustParse(NftsirenLogoData)
	LooksrareLogo = images.MustParse(LooksrareLogoData)
	MagicedenLogo = images.MustParse(MagicedenLogoData)
	OpenseaLogo   = images.MustParse(OpenseaLogoData)
	SolanaLogo    = images.MustParse(SolanaLogoData)
	TwitterLogo   = images.MustParse(TwitterLogoData)
)
