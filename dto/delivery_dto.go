package dto

type DeliveryOption struct {
	CourierName string `json:"courier_name"`
	CourierCode string `json:"courier_code"`
	CourierLogo string `json:"courier_logo"`
}

type DeliveryOptionUserMerchantResDTO struct {
	CourierName string `json:"courier_name"`
	CourierCode string `json:"courier_code"`
	CourierLogo string `json:"courier_logo"`
	IsChecked   bool   `json:"is_checked"`
}

type DeliveryGetAllOptionResDTO struct {
	DeliveryOptions []DeliveryOption `json:"delivery_options"`
	Total           int              `json:"total"`
}

type DeliveryGetMerchantOptionResDTO struct {
	MerchantDomain  string           `json:"merchant_domain"`
	MerchantName    string           `json:"merchant_name"`
	DeliveryOptions []DeliveryOption `json:"delivery_options"`
	Total           int              `json:"total"`
}

type DeliveryUpdateMerchantOptionReqDTO struct {
	CourierCode string `json:"courier_code"`
	IsChecked   bool   `json:"is_checked"`
}

type DeliveryUpdateMerchantOptionResDTO struct {
	CourierCode string `json:"courier_code"`
	IsChecked   bool   `json:"is_checked"`
}

type RajaOngkirDeliveryInfoReqDTO struct {
	Origin      int    `json:"origin"`
	Destination int    `json:"destination"`
	Weight      int    `json:"weight"`
	Courier     string `json:"courier"`
}

type RajaOngkirDeliveryInfoResDTO struct {
	Rajaongkir struct {
		Results []struct {
			Code  string `json:"code"`
			Name  string `json:"name"`
			Costs []struct {
				Service     string `json:"service"`
				Description string `json:"description"`
				Cost        []struct {
					Value float64 `json:"value"`
					Etd   string  `json:"etd"`
					Note  string  `json:"note"`
				} `json:"cost"`
			} `json:"costs"`
		} `json:"results"`
	} `json:"rajaongkir"`
}
