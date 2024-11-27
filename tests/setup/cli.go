package setup

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/types"
)

func CLI() (cli.Config, cli.CLI, types.KeyInfo, error) {
	cliAddress, err := types.NewAddressFromString("f1zrcoxi44hlmsdzfpwoce4hucmzqd5i4bwilg2ii")
	if err != nil {
		return cli.Config{}, cli.CLI{}, types.KeyInfo{}, fmt.Errorf("failed to generate filecoin address from string: %w", err)
	}

	encriptedPK := "7b2254797065223a22736563703235366b31222c22507269766174654b6579223a224a476741526d54516653764b2b38587a6e5a4b49703077547161574c516a476a4f69765767564856686d553d227d"

	cfg := cli.Config{
		Env: "testing",
		Route: cli.Route{
			Balance:   "/api/v1/balance",
			Banks:     "/api/v1/banks",
			Deposit:   "/api/v1/deposit",
			Withdraw:  "/api/v1/withdraw",
			Authorize: "/api/v1/authorize",
			Refund:    "/api/v1/refund",
		},
		Wallet: types.Wallet{
			Address: cliAddress,
		},
	}

	cl := cli.CLI{Validate: validator.New()}

	if err := cli.RegisterValidators(cl); err != nil {
		return cli.Config{}, cli.CLI{}, types.KeyInfo{}, fmt.Errorf("failed to register validators: %w", err)
	}

	pkOut := make([]byte, hex.DecodedLen(len(encriptedPK)))
	if _, err := hex.Decode(pkOut, []byte(encriptedPK)); err != nil {
		return cli.Config{}, cli.CLI{}, types.KeyInfo{}, fmt.Errorf("failed to decode private key: %w", err)
	}

	var keyInfo types.KeyInfo

	if err := json.Unmarshal(pkOut, &keyInfo); err != nil {
		return cli.Config{}, cli.CLI{}, types.KeyInfo{}, fmt.Errorf("failed to convert private key: %w", err)
	}

	return cfg, cl, keyInfo, nil
}
