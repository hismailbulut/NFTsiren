package nft

import "nftsiren/pkg/number"

func WeiToEth(wei number.Number) number.Number {
	return wei.Div(number.NewFromInt(1e18))
}

func EthToWei(eth number.Number) number.Number {
	return eth.Mul(number.NewFromInt(1e18))
}

func LamportsToSol(lamports number.Number) number.Number {
	return lamports.Div(number.NewFromInt(1e9))
}

func SolToLamports(sol number.Number) number.Number {
	return sol.Mul(number.NewFromInt(1e9))
}
