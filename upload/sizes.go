package upload

import "math"

func Auto() *int {
    return nil
}

func Px(px int) *int {
    return &px
}

func Rem(rem float64) *int {
    px := int(math.Ceil(rem * 16))

    return &px
}
