package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
	"github.com/mr-tron/base58"
)

func main() {
	compressFile()
	data := "N5mMy/LQRYthXLzGsaNnxHSen+9zBmIuGxtYkQEgvJoVnjerLY8vMAJs7kVxEuv8CaE7TRtQDy1ssBGokmhzKz9CDwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	// compress
	rawData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("raw data length", len(rawData))

	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		panic(err)
	}

	encodeRawData := encoder.EncodeAll(rawData, nil)
	fmt.Println("compressed data length: ", len(encodeRawData))
	// to base64
	fmt.Println(base64.StdEncoding.EncodeToString(encodeRawData))

	// zstd + base64
	dataEncoded := "KLUv/QBYhQIAZAQ3mYzL8tBFi2FcvMaxo2fEdJ6f73MGYi4bG1iRASC8mhWeN6stjy8wAmzuRXES6/wJoTtNG1APLWywEaiSaHMrP0IPAAEAAgAEJxJLcAY="
	content, err := base64.StdEncoding.DecodeString(dataEncoded)
	if err != nil {
		panic(err)
	}

	decoder, err := zstd.NewReader(nil)
	if err != nil {
		panic(err)
	}
	defer decoder.Close()

	res, err := decoder.DecodeAll(content, nil)
	if err != nil {
		panic(err)
	}

	// 	pub struct Account {
	//     /// The mint associated with this account
	//     pub mint: Pubkey,
	//     /// The owner of this account.
	//     pub owner: Pubkey,
	//     /// The amount of tokens this account holds.
	//     pub amount: u64,
	//     /// If `delegate` is `Some` then `delegated_amount` represents
	//     /// the amount authorized by the delegate
	//     pub delegate: COption<Pubkey>,
	//     /// The account's state
	//     pub state: AccountState,
	//     /// If is_some, this is a native token, and the value logs the rent-exempt reserve. An Account
	//     /// is required to be rent-exempt, so the value is used by the Processor to ensure that wrapped
	//     /// SOL accounts do not drop below this threshold.
	//     pub is_native: COption<u64>,
	//     /// The amount delegated
	//     pub delegated_amount: u64,
	//     /// Optional authority to close the account.
	//     pub close_authority: COption<Pubkey>,
	// }
	mint := base58.Encode(res[:32])
	owner := base58.Encode(res[32:64])
	amount := binary.LittleEndian.Uint64(res[64:72])
	delegate := base58.Encode(res[72:104])
	state := res[104:105][0]

	fmt.Printf("mint: %s, owner: %s, amount:%d delegate: %s state: %d\n", mint, owner, amount, delegate, state)
}

func compressFile() {
	in, err := os.OpenFile("../go.sum", os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer in.Close()
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile("go.sum.zstd", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	defer out.Close()
	w, err := zstd.NewWriter(out, zstd.WithEncoderLevel(zstd.EncoderLevel(3)))
	if err != nil {
		panic(err)
	}
	w.Close()
	i, err := io.Copy(w, in)
	if err != nil {
		panic(err)
	}
	w.Flush()
	fmt.Println(i)
}
