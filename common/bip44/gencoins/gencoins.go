package main

import (
	"strings"
)

const coins = `	TypeBitcoin               uint32 = 0x80000000
	TypeTestnet               uint32 = 0x80000001
	TypeLitecoin              uint32 = 0x80000002
	TypeDogecoin              uint32 = 0x80000003
	TypeReddcoin              uint32 = 0x80000004
	TypeDash                  uint32 = 0x80000005
	TypePeercoin              uint32 = 0x80000006
	TypeNamecoin              uint32 = 0x80000007
	TypeFeathercoin           uint32 = 0x80000008
	TypeCounterparty          uint32 = 0x80000009
	TypeBlackcoin             uint32 = 0x8000000a
	TypeNuShares              uint32 = 0x8000000b
	TypeNuBits                uint32 = 0x8000000c
	TypeMazacoin              uint32 = 0x8000000d
	TypeViacoin               uint32 = 0x8000000e
	TypeClearingHouse         uint32 = 0x8000000f
	TypeRubycoin              uint32 = 0x80000010
	TypeGroestlcoin           uint32 = 0x80000011
	TypeDigitalcoin           uint32 = 0x80000012
	TypeCannacoin             uint32 = 0x80000013
	TypeDigiByte              uint32 = 0x80000014
	TypeOpenAssets            uint32 = 0x80000015
	TypeMonacoin              uint32 = 0x80000016
	TypeClams                 uint32 = 0x80000017
	TypePrimecoin             uint32 = 0x80000018
	TypeNeoscoin              uint32 = 0x80000019
	TypeJumbucks              uint32 = 0x8000001a
	TypeziftrCOIN             uint32 = 0x8000001b
	TypeVertcoin              uint32 = 0x8000001c
	TypeNXT                   uint32 = 0x8000001d
	TypeBurst                 uint32 = 0x8000001e
	TypeMonetaryUnit          uint32 = 0x8000001f
	TypeZoom                  uint32 = 0x80000020
	TypeVpncoin               uint32 = 0x80000021
	TypeCanadaeCoin           uint32 = 0x80000022
	TypeShadowCash            uint32 = 0x80000023
	TypeParkByte              uint32 = 0x80000024
	TypePandacoin             uint32 = 0x80000025
	TypeStartCOIN             uint32 = 0x80000026
	TypeMOIN                  uint32 = 0x80000027
	TypeArgentum              uint32 = 0x8000002D
	TypeGlobalCurrencyReserve uint32 = 0x80000031
	TypeNovacoin              uint32 = 0x80000032
	TypeAsiacoin              uint32 = 0x80000033
	TypeBitcoindark           uint32 = 0x80000034
	TypeDopecoin              uint32 = 0x80000035
	TypeTemplecoin            uint32 = 0x80000036
	TypeAIB                   uint32 = 0x80000037
	TypeEDRCoin               uint32 = 0x80000038
	TypeSyscoin               uint32 = 0x80000039
	TypeSolarcoin             uint32 = 0x8000003a
	TypeSmileycoin            uint32 = 0x8000003b
	TypeEther                 uint32 = 0x8000003c
	TypeEtherClassic          uint32 = 0x8000003d
	TypeOpenChain             uint32 = 0x80000040
	TypeOKCash                uint32 = 0x80000045
	TypeDogecoinDark          uint32 = 0x8000004d
	TypeElectronicGulden      uint32 = 0x8000004e
	TypeClubCoin              uint32 = 0x8000004f
	TypeRichCoin              uint32 = 0x80000050
	TypePotcoin               uint32 = 0x80000051
	TypeQuarkcoin             uint32 = 0x80000052
	TypeTerracoin             uint32 = 0x80000053
	TypeGridcoin              uint32 = 0x80000054
	TypeAuroracoin            uint32 = 0x80000055
	TypeIXCoin                uint32 = 0x80000056
	TypeGulden                uint32 = 0x80000057
	TypeBitBean               uint32 = 0x80000058
	TypeBata                  uint32 = 0x80000059
	TypeMyriadcoin            uint32 = 0x8000005a
	TypeBitSend               uint32 = 0x8000005b
	TypeUnobtanium            uint32 = 0x8000005c
	TypeMasterTrader          uint32 = 0x8000005d
	TypeGoldBlocks            uint32 = 0x8000005e
	TypeSaham                 uint32 = 0x8000005f
	TypeChronos               uint32 = 0x80000060
	TypeUbiquoin              uint32 = 0x80000061
	TypeEvotion               uint32 = 0x80000062
	TypeSaveTheOcean          uint32 = 0x80000063
	TypeBigUp                 uint32 = 0x80000064
	TypeGameCredits           uint32 = 0x80000065
	TypeDollarcoins           uint32 = 0x80000066
	TypeZayedcoin             uint32 = 0x80000067
	TypeDubaicoin             uint32 = 0x80000068
	TypeStratis               uint32 = 0x80000069
	TypeShilling              uint32 = 0x8000006a
	TypePiggyCoin             uint32 = 0x80000076
	TypeMonero                uint32 = 0x80000080
	TypeNavCoin               uint32 = 0x80000082
	TypeFactomFactoids        uint32 = 0x80000083
	TypeFactomEntryCredits    uint32 = 0x80000084
	TypeZcash                 uint32 = 0x80000085
	TypeLisk                  uint32 = 0x80000086
	TypeFactomIdentity        uint32 = 0x80000119`

func main() {
	rows := strings.Split(coins, "\n")
	var coins []struct {
		coinName  string
		coinValue string
	}
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if row != "" {
			rowParts := strings.Split(row, "uint32")
			coinName := strings.TrimPrefix(strings.TrimSpace(rowParts[0]), "Type")
			coinValue := strings.TrimSpace(strings.Replace(rowParts[1], "=", "", -1))
			coins = append(coins, struct {
				coinName  string
				coinValue string
			}{
				coinName:  coinName,
				coinValue: coinValue,
			})
		}
	}
	for _, coin := range coins {
		println("coinTypes[`"+coin.coinName+"`]=", coin.coinValue)
	}
}
