// wei를 eth 단위로 변환
func WeiToEth(wei *big.Int) *big.Float {
	eth := new(big.Float)
	eth.SetString(wei.String())
	ethValue := new(big.Float).Quo(eth, big.NewFloat(math.Pow10(18)))
	return ethValue
}
