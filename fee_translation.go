package finnomena

import "github.com/jwitmann/finnomena-models"

var FeeDescriptionTranslation = map[string]string{
	"ค่าใช้จ่ายอื่นๆ":                                        "other fee",
	"ค่าธรรมเนียมการขายหน่วยลงทุน (Front-end Fee)":           "purchase fee",
	"ค่าธรรมเนียมการจัดการ":                                  "management fee",
	"ค่าธรรมเนียมการรับซื้อคืนหน่วยลงทุน (Back-end Fee)":     "redemption fee",
	"ค่าธรรมเนียมการสับเปลี่ยนหน่วยลงทุนเข้า (SWITCHING IN)": "switch in fee",
	"ค่าธรรมเนียมการสับเปลี่ยนหน่วยลงทุนออก (SWITCHING OUT)": "switch out fee",
	"ค่าธรรมเนียมและค่าใช้จ่ายรวมทั้งหมด":                    "total expense ratio",
	"ค่าธรรมเนียมการโอนหน่วยลงทุน":                           "unit transfer fee",
	"ค่าธรรมเนียมนายทะเบียนหน่วย":                            "registrar fee",
	"ค่าธรรมเนียมผู้ดูแลผลประโยชน์":                          "trustee fee",
}

var FeeOtherDescTranslation = map[string]string{
	"บริษัทจัดการอาจเรียกเก็บค่าธรรมเนียมการขายและค่าธรรมเนียมการรับซื้อคืนกับผู้ลงทุนแต่ละกลุ่มไม่เท่ากัน ทั้งนี้ ผู้ลงทุนสามารถดูรายละเอียดเพิ่มเติมได้ที่หนังสือชี้ชวนส่วนข้อมูลกองทุนรวม":                              "The fund management company may charge different sales fees and redemption fees to different investor groups. Investors can find more details in the fund information section of the prospectus.",
	"บริษัทจัดการอาจพิจารณาเปลี่ยนแปลงค่าธรรมเนียมที่เรียกเก็บจริงเพื่อให้สอดคล้องกับกลยุทธ์หรือค่าใช้จ่ายในการบริหารจัดการ และรวมค่าใช้จ่ายเป็นข้อมูลของรอบปีบัญชีล่าสุดหรือประมาณการเบื้องต้น (กรณียังไม่ครบรอบปีบัญชี)": "The management company may consider adjusting the actual fees charged to align with its strategy or management expenses, and include these expenses in the data for the latest accounting year or a preliminary estimate (if the accounting year has not yet ended).",
}

var FeeUnitTranslation = map[string]string{
	"ต่อปี ของมูลค่าทรัพย์สินสุทธิของกองทุน": "per year of NAV",
	"% ต่อปี": "% per year",
	"บาท":     "baht",
}

func TranslateFee(fee *models.Fee, useEnglish bool) {
	if !useEnglish {
		return
	}

	if trans, ok := FeeDescriptionTranslation[fee.Description]; ok {
		fee.Description = trans
	}

	if trans, ok := FeeOtherDescTranslation[fee.OtherDescription]; ok {
		fee.OtherDescription = trans
	}

	if trans, ok := FeeUnitTranslation[fee.Unit]; ok {
		fee.Unit = trans
	}
}
