package ordervalueobject

type OrderPrice struct {
	OriginPrice   int64
	DiscountPrice int64
	PaymentPrice  int64
}

func NewOrderPrice(originPrice, discountPrice, paymentPrice int64) (*OrderPrice, error) {
	return &OrderPrice{
		OriginPrice:   originPrice,
		DiscountPrice: discountPrice,
		PaymentPrice:  paymentPrice,
	}, nil
}
