import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func MakeKey() (string, string, error) {
	privateKey, err := ethCrypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := ethCrypto.FromECDSAPub(publicKeyECDSA)
	privateKeyBytes := ethCrypto.FromECDSA(privateKey)

	publicKeyHex := hexutil.Encode(publicKeyBytes)[4:]   // 0x와 04 제거
	privateKeyHex := hexutil.Encode(privateKeyBytes)[2:] // 0x 제거

	fmt.Printf("공개키: %s, 개인키: %s\n", publicKeyHex, privateKeyHex)

	return publicKeyHex, privateKeyHex, nil
}

func MakeAddress() (string, error) {
	privateKey, err := ethCrypto.GenerateKey()
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := ethCrypto.FromECDSAPub(publicKeyECDSA)
	privateKeyBytes := ethCrypto.FromECDSA(privateKey)

	publicKeyHex := hexutil.Encode(publicKeyBytes)[4:]   // 0x와 04 제거
	privateKeyHex := hexutil.Encode(privateKeyBytes)[2:] // 0x 제거

	fmt.Printf("공개키: %s, 개인키: %s\n", publicKeyHex, privateKeyHex)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	fmt.Printf("주소: %s\n", address)

	return address, nil
}

func createKeystore() {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "..."
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 지갑 주소
}

func importKeystore() {
	ks := keystore.NewKeyStore("./tmp2", keystore.StandardScryptN, keystore.StandardScryptP)

	filePath := "./tmp/UTC--2022-10-05T06-19-14.753221975Z--75a1d322d5c5bd1b98ae8061517f036e750f4f7c"

	jsonBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	password := "..."

	// account 가져오기
	account, err := ks.Import(jsonBytes, password, password)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())

	// 기존 파일 삭제
	if err := os.Remove(filePath); err != nil {
		log.Fatal(err)
	}
}
