package serializers

import "github.com/tuvistavie/structomap"

type PaymentPlanSerializer struct {
	*structomap.Base
}

func NewPaymentPlanSerializer() *PaymentPlanSerializer {

	serializer := &PaymentPlanSerializer{structomap.New()}
	serializer.UseCamelCase()
	serializer.Pick("Uuid", "Name", "Bandwidth", "Price", "Devices", "Duration")

	return serializer
}
