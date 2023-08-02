package model

type Body struct {
	Content string `json:"content"`
}

type UpsertProductPayload struct {
	ActivityID string `json:"id"`
	ClientKey  string `json:"client_key"`
	// Products   []Product `json:"products"`
}

type Product struct {
	BrandID         string           `json:"brand_id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	ImageURLs       []string         `json:"image_urls"`
	DimensionUnit   string           `json:"dimension_unit"`
	Width           float64          `json:"width"`
	Height          float64          `json:"height"`
	Length          float64          `json:"length"`
	WeightUnit      string           `json:"weight_unit"`
	Weight          float64          `json:"weight"`
	ChannelProducts []ChannelProduct `json:"channel_products"`
	ParentSKU       string           `json:"parent_sku"`
	HasVariant      bool             `json:"has_variant"`
}

type ChannelProduct struct {
	Channel              string           `json:"channel"`
	ChannelCategoryID    string           `json:"channel_category_id"`
	Attributes           []Attribute      `json:"attributes"`
	ChannelVariants      []ChannelVariant `json:"channel_variants"`
	ChannelProductStatus string           `json:"channel_product_status"`
}

type ChannelVariant struct {
	Name       string      `json:"variant_name"`
	Active     bool        `json:"active"`
	SKU        string      `json:"sku"`
	Stock      int         `json:"stock"`
	Price      float64     `json:"price"`
	IsPrimary  bool        `json:"is_primary,omitempty"`
	Attributes []Attribute `json:"attributes"`
	// RemoteVariant          json.RawMessage `json:"remote_variant"`
	ImageURLs              []string `json:"image_urls"`
	VariantRemoteWarehouse []string `json:"variant_remote_warehouse"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
