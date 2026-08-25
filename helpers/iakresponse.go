package helpers

type mainResponse struct {
	ProviderMsg string
	MainMsg     string
	Maincode    string
}

var respInq = map[string]mainResponse{
	// SUCCESS
	"00": {"INQUIRY SUCCESS", "Inquiry accepted. Product and pricing verified. Awaiting payment confirmation from merchant.", CodeInqSuccess},

	// ERR User Input / Data
	"01":  {"INVOICE HAS BEEN PAID", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransactionNoOrStatus},
	"02":  {"BILL UNPAID", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransactionNoOrStatus},
	"04":  {"BILLING ID EXPIRED", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidCustID},
	"06":  {"INQUIRY ID NOT FOUND", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidCustID},
	"08":  {"BILLING ID BLOCKED", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidIdGame},
	"09":  {"INQUIRY FAILED", "The execution request was rejected by the distribution core gateway due to structural rule conflicts.", CodeServiceDisruption},
	"10":  {"BILL IS NOT AVAILABLE", "The requested target product code or SKU is missing or inactive in our catalog.", CodeInvalidCustID},
	"42":  {"PAYMENT REQUEST HAVEN'T BEEN RECEIVED", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransactionNoOrStatus},
	"44":  {"EXCEEDING MAXIMAL DAILY INQUIRY ALLOWED", "The total value or volume of transactions has exceeded the agreed maximum daily limit.", CodeMaxTrx},
	"45":  {"TOO MANY INQUIRY REQUESTS", "The total value or volume of transactions has exceeded the agreed maximum daily limit.", CodeMaxTrx},
	"141": {"INVALID USER ID / ZONE ID / SERVER ID / ROLENAME", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidCustID},
	"142": {"INVALID USER ID", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidCustID},
	"143": {"INQUIRY NOT NEEDED", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransactionNoOrStatus},
	//belum kebaWAh====================
	// System / Merchant Error
	"91":  {"DATABASE CONNECTION ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeErrPvd4305},
	"92":  {"GENERAL ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeErrPvd4305},
	"93":  {"INVALID AMOUNT", "The destination identity string failed regex or basic character formatting check.", CodeInvalidCustId},
	"94":  {"SERVICE HAS EXPIRED", "The requested target product code or SKU is missing or inactive in our catalog.", CodeInvalidProductNotFound},
	"100": {"INVALID SIGNATURE", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeErrInt200},
	"101": {"INVALID COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeErrInt203},
	"102": {"INVALID IP ADDRESS", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeErrInt200},
	"103": {"TIMEOUT", "Upstream server failed to respond within the designated execution window. Money safely refunded.", CodeErrPvd1302},
	"105": {"MISC ERROR / BILLER SYSTEM ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeErrPvd4305},
	"106": {"PRODUCT IS TEMPORARILY OUT OF SERVICE", "The upstream vendor gateway for this specific product is currently undergoing maintenance.", CodeErrPvd1300},
	"107": {"XML / JSON FORMAT ERROR", "The destination identity string failed regex or basic character formatting check.", CodeInvalidCustId},
	"110": {"SYSTEM UNDER MAINTENANCE", "The upstream vendor gateway for this specific product is currently undergoing maintenance.", CodeErrPvd1300},
	"117": {"PAGE NOT FOUND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeErrInt203},
	"204": {"WRONG AUTHENTICATION", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeErrInt200},
	"205": {"WRONG COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeErrInt203},
}

var respPay = map[string]mainResponse{
	// SUCCESS
	"00": {"SUCCESS / PAYMENT SUCCESS", "Payment confirmed and transaction successfully processed. Delivery is in progress.", CodeSuccess},

	// PENDING
	"05":  {"UNDEFINED ERROR", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"39":  {"PENDING / TRANSACTION IN PROCESS", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"91":  {"DATABASE CONNECTION ERROR", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"94":  {"SERVICE HAS EXPIRED", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"103": {"TIMEOUT", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"105": {"MISC ERROR / BILLER SYSTEM ERROR", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"110": {"SYSTEM UNDER MAINTENANCE", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"201": {"UNDEFINED RESPONSE CODE", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},

	// FAILED
	"92":  {"GENERAL ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeServiceDisruption},
	"93":  {"INVALID AMOUNT", "The destination identity string failed regex or basic character formatting check.", CodeServiceDisruption},
	"100": {"INVALID SIGNATURE", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeServiceDisruption},
	"101": {"INVALID COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeServiceDisruption},
	"102": {"INVALID IP ADDRESS", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeServiceDisruption},
	"106": {"PRODUCT IS TEMPORARILY OUT OF SERVICE", "The upstream vendor gateway for this specific product is currently undergoing maintenance.", CodeServiceDisruption},
	"107": {"XML / JSON FORMAT ERROR", "The destination identity string failed regex or basic character formatting check.", CodeServiceDisruption},
	"117": {"PAGE NOT FOUND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeServiceDisruption},
	"204": {"WRONG AUTHENTICATION", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeServiceDisruption},
	"205": {"WRONG COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeServiceDisruption},
}

// ConvertIAKPaymentResponse converts IAK payment response code (rc) to main response format.
func ConvertIAKPaymentResponse(rc string) (mainCode string, providerMsg string, mainMsg string) {
	if resp, exists := respPay[rc]; exists {
		return resp.Maincode, resp.ProviderMsg, resp.MainMsg
	}
	return CodeErrPvd3303, "PAYMENT FAILED", "Upstream provider payment failed."
}
