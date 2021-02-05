package faucet

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func (f *Faucet) Send(recipient string) error {
	for _, denom := range f.denoms {
		if _, err := f.cliexec([]string{"tx", "bank", "send", f.keyName, recipient,
			fmt.Sprintf("%d%s", f.creditAmount, denom), "--yes", "--chain-id", f.chainID,
			"--keyring-backend", "test"},
			f.keyringPassword, f.keyringPassword, f.keyringPassword,
		); err != nil {
			return err
		}
	}

	return nil
}

func (f *Faucet) GetTotalSent(recipient string) (uint64, error) {
	args := []string{
		"query", "txs", "--events",
		fmt.Sprintf("message.sender=%s&transfer.recipient=%s", f.faucetAddress, recipient),
		"--page", "1",
		"--limit", "1000",
		"--keyring-backend", "test",
		// "--trust-node",
	}

	output, err := f.cliexec(args)
	if err != nil {
		return 0, err
	}

	var result types.SearchTxsResult
	if err := f.cdc.UnmarshalJSON([]byte(output), &result); err != nil {
		fmt.Println(">>>>>", output)
		return 0, err
	}

	var total uint64

	for _, tx := range result.Txs {
		stdTx := tx.Tx.(auth.StdTx)
		if len(stdTx.Msgs) == 0 {
			return 0, fmt.Errorf("no MsgSend available in transaction")
		}

		msg := stdTx.Msgs[0].(bank.MsgSend)
		for _, denom := range f.denoms {
			total += msg.Amount.AmountOf(denom).Uint64()
		}
	}

	return total, nil
}
