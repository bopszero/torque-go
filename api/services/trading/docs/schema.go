package docs

var doc = "" +
	"openapi: 3.0.0\n" +
	"info:\n" +
	"  version: 1.0.0\n" +
	"  title: Torque Wallet Service\n" +
	"servers:\n" +
	"  - url: \"/\"\n" +
	"components:\n" +
	"  securitySchemes:\n" +
	"    BearerAuth:\n" +
	"      type: http\n" +
	"      scheme: bearer\n" +
	"  schemas:\n" +
	"    Paging:\n" +
	"      type: object\n" +
	"      properties:\n" +
	"        limit:\n" +
	"          type: integer\n" +
	"          example: 10\n" +
	"        offset:\n" +
	"          type: integer\n" +
	"          example: 0\n" +
	"        before_id:\n" +
	"          type: integer\n" +
	"        after_id:\n" +
	"          type: integer\n" +
	"      required:\n" +
	"        - limit\n" +
	"    Decimal:\n" +
	"      type: string\n" +
	"      description: Decimal number in string\n" +
	"      example: \"14211.364041\"\n" +
	"    ResponseMessage:\n" +
	"      type: string\n" +
	"      example: \"Field validation for 'user_id' failed on the 'required' tag.\"\n" +
	"    ID:\n" +
	"      type: integer\n" +
	"      minimum: 1\n" +
	"      example: 8466032601\n" +
	"    Reference:\n" +
	"      type: string\n" +
	"      example: \"9127400572\"\n" +
	"    Currency:\n" +
	"      type: string\n" +
	"      enum:\n" +
	"        - BTC\n" +
	"        - BCH\n" +
	"        - LTC\n" +
	"        - ETH\n" +
	"        - TRX\n" +
	"        - USDT\n" +
	"        - TORQ\n" +
	"    Network:\n" +
	"      type: string\n" +
	"      enum:\n" +
	"        - BTC:TESTNET\n" +
	"        - BTC\n" +
	"        - BCH:TESTNET\n" +
	"        - BCH\n" +
	"        - LTC:TESTNET\n" +
	"        - LTC\n" +
	"        - ETH:TEST_ROPSTEN\n" +
	"        - ETH\n" +
	"        - TRX:TEST_SHASTA\n" +
	"        - TRX\n" +
	"        - TORQ\n" +
	"    Address:\n" +
	"      type: string\n" +
	"      minLength: 16\n" +
	"      maxLength: 255\n" +
	"      example: bc1q74m8sy7dpqwuegcrddfwqp3gzwau62k7rys0my\n" +
	"    Timestamp:\n" +
	"      type: integer\n" +
	"      example: 1578396200\n" +
	"    AuthCode:\n" +
	"      type: string\n" +
	"      minLength: 6\n" +
	"      maxLength: 6\n" +
	"      example: \"111111\"\n" +
	"    Note:\n" +
	"      type: string\n" +
	"      example: I lay my love on you...\n" +
	"    Token32:\n" +
	"      type: string\n" +
	"      example: \"b5bac6a21c124fbaa56b7688720cf407\"\n" +
	"security:\n" +
	"  - BearerAuth: [ ]\n" +
	"paths:\n" +
	"  /v1/txn/deposit/export/:\n" +
	"    post:\n" +
	"      summary: Generate Deposit export token\n" +
	"      tags:\n" +
	"        - Transaction\n" +
	"        - Deposit\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"  /v1/txn/deposit/export/{token}/:\n" +
	"    post:\n" +
	"      summary: Download Deposit report\n" +
	"      tags:\n" +
	"        - Transaction\n" +
	"        - Deposit\n" +
	"      parameters:\n" +
	"        - in: path\n" +
	"          name: token\n" +
	"          required: true\n" +
	"          schema:\n" +
	"            $ref: \"#/components/schemas/Token32\"\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: Exported file binary data\n" +
	"        \"500\":\n" +
	"          description: Server error\n" +
	"  /v1/txn/withdrawal/export/:\n" +
	"    post:\n" +
	"      summary: Generate Withdrawal export token\n" +
	"      tags:\n" +
	"        - Transaction\n" +
	"        - Withdrawal\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"  /v1/txn/withdrawal/export/{token}/:\n" +
	"    post:\n" +
	"      summary: Download Withdrawal report\n" +
	"      tags:\n" +
	"        - Transaction\n" +
	"        - Withdrawal\n" +
	"      parameters:\n" +
	"        - in: path\n" +
	"          name: token\n" +
	"          required: true\n" +
	"          schema:\n" +
	"            $ref: \"#/components/schemas/Token32\"\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: Exported file binary data\n" +
	"        \"500\":\n" +
	"          description: Server error\n" +
	""
