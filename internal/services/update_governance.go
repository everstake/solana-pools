package services

import (
	"fmt"
)

func (s Imp) UpdateGovernance() error {
	gov, err := s.DAO.GetGovernance(nil)
	if err != nil {
		return fmt.Errorf("UpdateGovernance: %w", err)
	}

	for _, governance := range gov {
		if governance.GeckoKey == "null" {
			continue
		}
		coin, err := s.coinGecko.CoinsID(governance.GeckoKey,
			false,
			false,
			true,
			false,
			false,
			false)

		if err != nil {
			return fmt.Errorf("UpdateGovernance: %w", err)
		}

		if coin.MarketData.TotalSupply != nil {
			governance.MaximumTokenSupply = *coin.MarketData.TotalSupply
		}

		governance.CirculatingSupply = coin.MarketData.CirculatingSupply

		governance.USD = coin.MarketData.CurrentPrice["usd"]

		governance.Image = coin.Image.Large
	}

	if err := s.DAO.SaveGovernance(gov...); err != nil {
		return fmt.Errorf("UpdateGovernance: %w", err)
	}

	return nil
}
