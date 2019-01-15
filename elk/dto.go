package elk

type BulkHeader struct {
	Index struct {
		Index string `json:"_index"`
		Type  string `json:"_type"`
	} `json:"index"`
}

type BulkMessage struct {
	Timestamp string `json:"@timestamp"`
	C         struct {
		Firm string `json:"firm"`
		Msg  string `json:"msg"`
		Name string `json:"name"`
	} `json:"c"`
	Message struct {
		MiFIDShortcodes struct {
			ClientIdentificationShortcode    int64 `json:"clientIdentificationShortcode"`
			InvestmentDecisionWFirmShortCode int64 `json:"investmentDecisionWFirmShortCode"`
			NonExecutingBrokerShortCode      int64 `json:"nonExecutingBrokerShortCode"`
		} `json:"MiFIDShortcodes"`
		AccountType                  int64         `json:"accountType"`
		AccountTypeName              string        `json:"accountType_name"`
		ClMsgSeqNum                  int64         `json:"clMsgSeqNum"`
		ClientOrderID                int64         `json:"clientOrderID"`
		DarkExecutionInstruction     []interface{} `json:"darkExecutionInstruction"`
		EMM                          int64         `json:"eMM"`
		EMMName                      string        `json:"eMM_name"`
		ExecutionInstruction         []interface{} `json:"executionInstruction"`
		ExecutionWithinFirmShortCode int64         `json:"executionWithinFirmShortCode"`
		FirmID                       string        `json:"firmID"`
		LPRole                       int64         `json:"lPRole"`
		LPRoleName                   string        `json:"lPRole_name"`
		MiFIDIndicators              []interface{} `json:"miFIDIndicators"`
		OrderPx                      int64         `json:"orderPx"`
		OrderQty                     string        `json:"orderQty"`
		OrderSide                    int64         `json:"orderSide"`
		OrderSideName                string        `json:"orderSide_name"`
		OrderType                    int64         `json:"orderType"`
		OrderTypeName                string        `json:"orderType_name"`
		STPID                        int64         `json:"sTPID"`
		SendingTime                  string        `json:"sendingTime"`
		SymbolIndex                  int64         `json:"symbolIndex"`
		TimeInForce                  int64         `json:"timeInForce"`
		TimeInForceName              string        `json:"timeInForce_name"`
		TradingCapacity              int64         `json:"tradingCapacity"`
		TradingCapacityName          string        `json:"tradingCapacity_name"`
	} `json:"message"`
	Raw string `json:"raw"`
	S   struct {
		ClOrdIDPrefix   string `json:"ClOrdIDPrefix"`
		EnteringFirmID  string `json:"EnteringFirmID"`
		ExecutingFirmID string `json:"ExecutingFirmID"`
		FirmName        string `json:"FirmName"`
		LogicalAccessID string `json:"LogicalAccessID"`
		OptiqSegment    string `json:"OptiqSegment"`
		PartitionID     string `json:"PartitionID"`
		PartitionName   string `json:"PartitionName"`
		PartitionType   string `json:"PartitionType"`
		Protocol        string `json:"Protocol"`
	} `json:"s"`
	Sbe string `json:"sbe"`
	T   struct {
		Cap     int64  `json:"cap"`
		CapTxt  string `json:"cap_txt"`
		DOegv   int64  `json:"d_oegv"`
		OegOut  string `json:"oeg_out"`
		OegvOut string `json:"oegv_out"`
	} `json:"t"`
	Tcp struct {
		Client struct {
			IP   string `json:"ip"`
			Port string `json:"port"`
		} `json:"client"`
		ID  string `json:"id"`
		Oeg struct {
			IP   string `json:"ip"`
			Port string `json:"port"`
		} `json:"oeg"`
		Sender string `json:"sender"`
	} `json:"tcp"`
	V []string `json:"v"`
}
