package azure

type sttResponse struct {
	RecognitionStatus string `json:"RecognitionStatus"`
	Offset            int    `json:"Offset"`
	Duration          int    `json:"Duration"`
	DisplayText       string `json:"DisplayText"`
}
