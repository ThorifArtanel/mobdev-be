package common

import (
	"os"
	"strconv"
)

func GetObjectURL() string {
	bucket := os.Getenv("OBJECT_URL")

	return bucket
}

func GetTokenExp() int {
	exp, _ := strconv.Atoi(os.Getenv("JWT_EXPIRED_HOUR"))
	return exp
}

func GetTokenSecret() string {
	return os.Getenv("SECRET")
}

func GetDBURL() string {
	return os.Getenv("DATABASE_URL")
}

func GetDescTruncLen() int {
	return 25
}

func DokterRole() string {
	return "DOKTER"
}

func KeluargaRole() string {
	return "KELUARGA"
}
