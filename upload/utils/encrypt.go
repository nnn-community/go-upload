package utils

import (
    "crypto/md5"
    "encoding/hex"
    "strconv"
    "time"
)

func GetMD5Hash(text string) string {
    text = strconv.Itoa(int(time.Now().Unix())) + "-" + text
    hash := md5.Sum([]byte(text))

    return hex.EncodeToString(hash[:])
}
