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
	"    DirectionType:\n" +
	"      type: integer\n" +
	"      enum:\n" +
	"        - -1\n" +
	"        - +1\n" +
	"      description: |\n" +
	"        * -1 - Outcome\n" +
	"        * +1 - Income\n" +
	"    ChannelType:\n" +
	"      type: integer\n" +
	"      enum:\n" +
	"        - 10000\n" +
	"        - 10001\n" +
	"        - 10002\n" +
	"        - 10003\n" +
	"        - 10004\n" +
	"        - 10100\n" +
	"        - 10101\n" +
	"        - 100000\n" +
	"        - 100001\n" +
	"        - 100002\n" +
	"        - 100003\n" +
	"        - 100004\n" +
	"        - 100100\n" +
	"        - 100101\n" +
	"        - 110000\n" +
	"        - 110001\n" +
	"        - 110002\n" +
	"      description: |\n" +
	"        * 10000 - Source Balance\n" +
	"        * 10001 - Source System\n" +
	"        * 10002 - Source Transfer\n" +
	"        * 10003 - Source Blockchain Network\n" +
	"        * 10004 - Source Torque Purchase\n" +
	"        * 10100 - Source Promo Code\n" +
	"        * 10101 - Source Trading Reward\n" +
	"        *\n" +
	"        * 100000 - Destination Balance\n" +
	"        * 100001 - Destination System\n" +
	"        * 100002 - Destination Transfer\n" +
	"        * 100003 - Destination Blockchain Network\n" +
	"        * 100004 - Destination Torque Purchase\n" +
	"        * 100100 - Destination Profit Reinvest\n" +
	"        * 100101 - Destination Profit Withdraw\n" +
	"        * 110000 - Destination Merchant Gorilla Hotel\n" +
	"        * 110001 - Destination Merchant Gorilla Flight\n" +
	"        * 110002 - Destination Merchant Torque Mall\n" +
	"    OrderStatus:\n" +
	"      type: integer\n" +
	"      enum:\n" +
	"        - 1\n" +
	"        - 2\n" +
	"        - 3\n" +
	"        - 4\n" +
	"        - 50\n" +
	"        - 98\n" +
	"        - 99\n" +
	"        - 100\n" +
	"        - 101\n" +
	"        - -1\n" +
	"        - -2\n" +
	"        - -3\n" +
	"      description: |\n" +
	"        * -1 - Failed\n" +
	"        * -2 - Canceled\n" +
	"        * -3 - Expired\n" +
	"        *\n" +
	"        * 1 - New\n" +
	"        * 2 - Init\n" +
	"        * 3 - Handle Source\n" +
	"        * 4 - Handle Destination\n" +
	"        *\n" +
	"        * 50 - Need Staff\n" +
	"        *\n" +
	"        * 97 - Failing\n" +
	"        * 98 - Refunding\n" +
	"        * 99 - Completing\n" +
	"        * 100 - Completed\n" +
	"        * 101 - Refunded\n" +
	"      example: 100\n" +
	"    Order:\n" +
	"      type: object\n" +
	"      properties:\n" +
	"        id:\n" +
	"          $ref: \"#/components/schemas/ID\"\n" +
	"        uid:\n" +
	"          $ref: \"#/components/schemas/ID\"\n" +
	"        direction_type:\n" +
	"          $ref: \"#/components/schemas/DirectionType\"\n" +
	"        currency:\n" +
	"          $ref: \"#/components/schemas/Currency\"\n" +
	"        status:\n" +
	"          $ref: \"#/components/schemas/OrderStatus\"\n" +
	"        src_channel_type:\n" +
	"          $ref: \"#/components/schemas/ChannelType\"\n" +
	"        src_channel_id:\n" +
	"          type: integer\n" +
	"          description: Related object ID of Source channel\n" +
	"        src_channel_ref:\n" +
	"          $ref: \"#/components/schemas/Reference\"\n" +
	"        src_channel_amount:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        src_channel_context:\n" +
	"          $ref: \"#/components/schemas/OrderChannelContext\"\n" +
	"        dst_channel_type:\n" +
	"          $ref: \"#/components/schemas/ChannelType\"\n" +
	"        dst_channel_id:\n" +
	"          type: integer\n" +
	"          description: Related object ID of Desctination channel\n" +
	"        dst_channel_ref:\n" +
	"          $ref: \"#/components/schemas/Reference\"\n" +
	"        dst_channel_amount:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        dst_channel_context:\n" +
	"          $ref: \"#/components/schemas/OrderChannelContext\"\n" +
	"        amount_sub_total:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        amount_fee:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        amount_discount:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        amount_total:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        note:\n" +
	"          $ref: \"#/components/schemas/Note\"\n" +
	"        create_time:\n" +
	"          $ref: \"#/components/schemas/Timestamp\"\n" +
	"        update_time:\n" +
	"          $ref: \"#/components/schemas/Timestamp\"\n" +
	"      required:\n" +
	"        - id\n" +
	"        - uid\n" +
	"        - direction_type\n" +
	"        - currency\n" +
	"        - status\n" +
	"        - src_channel_type\n" +
	"        - src_channel_id\n" +
	"        - src_channel_ref\n" +
	"        - src_channel_amount\n" +
	"        - dst_channel_type\n" +
	"        - dst_channel_id\n" +
	"        - dst_channel_ref\n" +
	"        - dst_channel_amount\n" +
	"        - amount_sub_total\n" +
	"        - amount_fee\n" +
	"        - amount_discount\n" +
	"        - amount_total\n" +
	"        - note\n" +
	"        - create_time\n" +
	"        - update_time\n" +
	"    OrderChannelContext:\n" +
	"      type: object\n" +
	"      properties:\n" +
	"        meta:\n" +
	"          type: object\n" +
	"          nullable: true\n" +
	"        details:\n" +
	"          type: object\n" +
	"          nullable: true\n" +
	"    LockAction:\n" +
	"      type: integer\n" +
	"      enum:\n" +
	"        - 1\n" +
	"        - 2\n" +
	"        - 3\n" +
	"        - 4\n" +
	"        - 5\n" +
	"      description: |\n" +
	"        * 1 - Personal-Send\n" +
	"        * 2 - Personal-TORQ-Transfer\n" +
	"        * 3 - Personal-TORQ-Reallocate\n" +
	"        * 4 - Personal-TORQ-Purchase\n" +
	"        * 5 - Trading-Send(Withdraw)\n" +
	"    CurrencyInfo:\n" +
	"      type: object\n" +
	"      properties:\n" +
	"        currency:\n" +
	"          $ref: \"#/components/schemas/Currency\"\n" +
	"        price_usd:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        price_usdt:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        priority:\n" +
	"          type: integer\n" +
	"          example: 11\n" +
	"        is_fiat:\n" +
	"          type: boolean\n" +
	"        decimal_places:\n" +
	"          type: integer\n" +
	"          example: 8\n" +
	"        icon_url:\n" +
	"          type: string\n" +
	"          example: https://torquebot.net/assests/coins/btc.png\n" +
	"        banner_url:\n" +
	"          type: string\n" +
	"          example: https://qc.torquebot.net/assests/coins/banner_btc.jpg\n" +
	"        symbol:\n" +
	"          type: string\n" +
	"          example: \"₿\"\n" +
	"        color_hex:\n" +
	"          type: string\n" +
	"          example: f8932b\n" +
	"        status_message:\n" +
	"          type: string\n" +
	"          example: \"Will be available on 2021-01-01\"\n" +
	"security:\n" +
	"  - BearerAuth: [ ]\n" +
	"tags:\n" +
	"  - name: Portfolio\n" +
	"  - name: Payment\n" +
	"paths:\n" +
	"  /v1/meta/handshake/:\n" +
	"    post:\n" +
	"      summary: Init meta data\n" +
	"      tags:\n" +
	"        - Meta\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                firebase_token:\n" +
	"                  type: string\n" +
	"                  example: fpnf2XHOS6O82VWajh8fCL:APA91bHc3SpnCwvPxHo0UlMdgEBVgzcJ_QWs80Y5-IPCynyt2FZR27wzr6Q86ZoY35kyUoYqZZlZRUFGjjhDlCIS6cMPwz7fXpy_hrzSX8ktkcgcsToMTtNfEcxz1V__wg5rr7xIIovu\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      features:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          type: object\n" +
	"                          properties:\n" +
	"                            code:\n" +
	"                              type: string\n" +
	"                              example: kyc\n" +
	"                            name:\n" +
	"                              type: string\n" +
	"                              example: KYC\n" +
	"                            is_available:\n" +
	"                              type: boolean\n" +
	"                      blockchain_networks:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          type: object\n" +
	"                          properties:\n" +
	"                            code:\n" +
	"                              $ref: \"#/components/schemas/Network\"\n" +
	"                            currency:\n" +
	"                              $ref: \"#/components/schemas/Currency\"\n" +
	"                            name:\n" +
	"                              type: string\n" +
	"                              example: Ethereum\n" +
	"                            token_transfer_code_name:\n" +
	"                              type: string\n" +
	"                              example: ERC20\n" +
	"                      network_currencies:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          type: object\n" +
	"                          properties:\n" +
	"                            currency:\n" +
	"                              $ref: \"#/components/schemas/Currency\"\n" +
	"                            network:\n" +
	"                              $ref: \"#/components/schemas/Network\"\n" +
	"                            priority:\n" +
	"                              type: integer\n" +
	"                              example: 2\n" +
	"                            withdrawal_fee:\n" +
	"                              $ref: \"#/components/schemas/Decimal\"\n" +
	"  /v1/meta/currency-info/:\n" +
	"    post:\n" +
	"      summary: Get system currency info\n" +
	"      tags:\n" +
	"        - Meta\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      currencies:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          $ref: \"#/components/schemas/CurrencyInfo\"\n" +
	"  /v1/portfolio/overview/get/:\n" +
	"    post:\n" +
	"      summary: Get Portfolio Overview\n" +
	"      tags:\n" +
	"        - Portfolio\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      balances:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          type: object\n" +
	"                          properties:\n" +
	"                            uid:\n" +
	"                              $ref: \"#/components/schemas/ID\"\n" +
	"                            is_available:\n" +
	"                              type: boolean\n" +
	"                            currency:\n" +
	"                              $ref: \"#/components/schemas/Currency\"\n" +
	"                            currency_price_usd:\n" +
	"                              $ref: \"#/components/schemas/Decimal\"\n" +
	"                            currency_price_usdt:\n" +
	"                              $ref: \"#/components/schemas/Decimal\"\n" +
	"                            amount:\n" +
	"                              $ref: \"#/components/schemas/Decimal\"\n" +
	"                            update_time:\n" +
	"                              $ref: \"#/components/schemas/Timestamp\"\n" +
	"  /v1/portfolio/currency/get/:\n" +
	"    post:\n" +
	"      summary: Get User Currency Info\n" +
	"      tags:\n" +
	"        - Portfolio\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      currency:\n" +
	"                        $ref: \"#/components/schemas/Currency\"\n" +
	"                      price_usd:\n" +
	"                        $ref: \"#/components/schemas/Decimal\"\n" +
	"                      price_usdt:\n" +
	"                        $ref: \"#/components/schemas/Decimal\"\n" +
	"                      account_no:\n" +
	"                        $ref: \"#/components/schemas/Address\"\n" +
	"                      balance:\n" +
	"                        $ref: \"#/components/schemas/Decimal\"\n" +
	"                      notice:\n" +
	"                        type: string\n" +
	"                        example: \"We cannot call APIs due to lacking of the sunlight in the midnight.\"\n" +
	"  /v1/portfolio/currency/order/list/:\n" +
	"    post:\n" +
	"      summary: List Orders of a Currency\n" +
	"      tags:\n" +
	"        - Portfolio\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                paging:\n" +
	"                  $ref: \"#/components/schemas/Paging\"\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      paging:\n" +
	"                        $ref: \"#/components/schemas/Paging\"\n" +
	"                      items:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          $ref: \"#/components/schemas/Order\"\n" +
	"  /v1/portfolio/currency/order/export/:\n" +
	"    post:\n" +
	"      summary: Get Token for exporting Orders of a Currency\n" +
	"      tags:\n" +
	"        - Portfolio\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                from_time:\n" +
	"                  $ref: \"#/components/schemas/Timestamp\"\n" +
	"                to_time:\n" +
	"                  $ref: \"#/components/schemas/Timestamp\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - from_time\n" +
	"                - to_time\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      token:\n" +
	"                        type: string\n" +
	"                        example: b5bac6a21c124fbaa56b7688720cf407\n" +
	"  /v1/portfolio/currency/order/export/{token}/:\n" +
	"    get:\n" +
	"      summary: Export Orders of a Currency\n" +
	"      tags:\n" +
	"        - Portfolio\n" +
	"      parameters:\n" +
	"        - in: path\n" +
	"          name: token\n" +
	"          required: true\n" +
	"          schema:\n" +
	"            type: string\n" +
	"            example: b5bac6a21c124fbaa56b7688720cf407\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: Exported file data\n" +
	"          content:\n" +
	"            text/csv:\n" +
	"              example: \"Code,Currency,Type,Status,Amount,Time,Note\"\n" +
	"        \"500\":\n" +
	"          description: Server error\n" +
	"  /v1/portfolio/currency/order/get/:\n" +
	"    post:\n" +
	"      summary: Get an Order\n" +
	"      tags:\n" +
	"        - Portfolio\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                ref:\n" +
	"                  $ref: \"#/components/schemas/Reference\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - ref\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    $ref: \"#/components/schemas/Order\"\n" +
	"  /v1/payment/order/checkout/:\n" +
	"    post:\n" +
	"      summary: Checkout a new Order\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                src_channel_type:\n" +
	"                  $ref: \"#/components/schemas/ChannelType\"\n" +
	"                src_channel_id:\n" +
	"                  type: integer\n" +
	"                src_channel_ref:\n" +
	"                  $ref: \"#/components/schemas/Reference\"\n" +
	"                src_channel_amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                src_channel_context:\n" +
	"                  type: object\n" +
	"                dst_channel_type:\n" +
	"                  $ref: \"#/components/schemas/ChannelType\"\n" +
	"                dst_channel_id:\n" +
	"                  type: integer\n" +
	"                dst_channel_ref:\n" +
	"                  $ref: \"#/components/schemas/Reference\"\n" +
	"                dst_channel_amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                dst_channel_context:\n" +
	"                  type: object\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - src_channel_type\n" +
	"                - src_channel_amount\n" +
	"                - dst_channel_type\n" +
	"                - dst_channel_amount\n" +
	"            examples:\n" +
	"              transfer_p2p:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"currency\": \"TORQ\",\n" +
	"                    \"src_channel_type\": 10000,\n" +
	"                    \"src_channel_amount\": \"12.34567898\",\n" +
	"                    \"dst_channel_type\": 100002,\n" +
	"                    \"dst_channel_amount\": \"12.34567898\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"user_identity\": \"powerman\",\n" +
	"                        \"note\": \"Hello World!\",\n" +
	"                      },\n" +
	"                  }\n" +
	"              profit_reinvest:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"currency\": \"TORQ\",\n" +
	"                    \"src_channel_type\": 10000,\n" +
	"                    \"src_channel_amount\": \"20\",\n" +
	"                    \"dst_channel_type\": 100100,\n" +
	"                    \"dst_channel_amount\": \"20\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"user_identity\": \"powerman\",\n" +
	"                        \"currency\": \"USDT\",\n" +
	"                        \"exchange_rate\": \"0.05\",\n" +
	"                      },\n" +
	"                  }\n" +
	"              profit_withdraw:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"currency\": \"TORQ\",\n" +
	"                    \"src_channel_type\": 10000,\n" +
	"                    \"src_channel_amount\": \"25\",\n" +
	"                    \"dst_channel_type\": 100101,\n" +
	"                    \"dst_channel_amount\": \"25\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"address\": \"0x83a7663b2b9d6d3f377a41d03b03ba0021e2f831\",\n" +
	"                        \"currency\": \"USDT\",\n" +
	"                        \"exchange_rate\": \"0.05\",\n" +
	"                      },\n" +
	"                  }\n" +
	"              blockchain_txn:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"currency\": \"USDT\",\n" +
	"                    \"src_channel_type\": 10003,\n" +
	"                    \"src_channel_amount\": \"1\",\n" +
	"                    \"dst_channel_type\": 100003,\n" +
	"                    \"dst_channel_amount\": \"1\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"to_address\": \"0xf12Db5BeC81DD6f41366EbcC599508212b751479\",\n" +
	"                      },\n" +
	"                  }\n" +
	"              torque_purchase:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"currency\": \"USDT\",\n" +
	"                    \"src_channel_type\": 10003,\n" +
	"                    \"src_channel_amount\": \"2\",\n" +
	"                    \"dst_channel_type\": 100004,\n" +
	"                    \"dst_channel_amount\": \"2\",\n" +
	"                  }\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      channel_src_info:\n" +
	"                        type: object\n" +
	"                        nullable: true\n" +
	"                      channel_dst_info:\n" +
	"                        type: object\n" +
	"                        nullable: true\n" +
	"                        example:\n" +
	"                          { \"fee\": { \"currency\": \"TORQ\", \"value\": \"0.001\" } }\n" +
	"  /v1/payment/order/init/:\n" +
	"    post:\n" +
	"      summary: Init a new Order\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                src_channel_type:\n" +
	"                  $ref: \"#/components/schemas/ChannelType\"\n" +
	"                src_channel_id:\n" +
	"                  type: integer\n" +
	"                src_channel_ref:\n" +
	"                  $ref: \"#/components/schemas/Reference\"\n" +
	"                src_channel_amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                src_channel_context:\n" +
	"                  type: object\n" +
	"                dst_channel_type:\n" +
	"                  $ref: \"#/components/schemas/ChannelType\"\n" +
	"                dst_channel_id:\n" +
	"                  type: integer\n" +
	"                dst_channel_ref:\n" +
	"                  $ref: \"#/components/schemas/Reference\"\n" +
	"                dst_channel_amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                dst_channel_context:\n" +
	"                  type: object\n" +
	"                amount_sub_total:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                amount_total:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                note:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - src_channel_type\n" +
	"                - src_channel_amount\n" +
	"                - dst_channel_type\n" +
	"                - dst_channel_amount\n" +
	"                - amount_sub_total\n" +
	"                - amount_total\n" +
	"            examples:\n" +
	"              transfer_p2p:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"currency\": \"TORQ\",\n" +
	"                    \"src_channel_type\": 10000,\n" +
	"                    \"src_channel_amount\": \"12.34567898\",\n" +
	"                    \"dst_channel_type\": 100002,\n" +
	"                    \"dst_channel_amount\": \"12.34567898\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"note\": \"Hello World!\",\n" +
	"                        \"user_identity\": \"powerman\",\n" +
	"                      },\n" +
	"                    \"amount_sub_total\": \"12.34567898\",\n" +
	"                    \"amount_total\": \"12.34567898\",\n" +
	"                    \"note\": \"Test transfer\",\n" +
	"                  }\n" +
	"              profit_reinvest:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"currency\": \"TORQ\",\n" +
	"                    \"src_channel_type\": 10000,\n" +
	"                    \"src_channel_amount\": \"20\",\n" +
	"                    \"dst_channel_type\": 100100,\n" +
	"                    \"dst_channel_amount\": \"20\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"user_identity\": \"powerman\",\n" +
	"                        \"currency\": \"USDT\",\n" +
	"                        \"exchange_rate\": \"0.05\",\n" +
	"                      },\n" +
	"                    \"amount_sub_total\": \"20\",\n" +
	"                    \"amount_total\": \"20\",\n" +
	"                    \"note\": \"Test reinvest USDT\",\n" +
	"                  }\n" +
	"              profit_withdraw:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"currency\": \"TORQ\",\n" +
	"                    \"src_channel_type\": 10000,\n" +
	"                    \"src_channel_amount\": \"25\",\n" +
	"                    \"dst_channel_type\": 100101,\n" +
	"                    \"dst_channel_amount\": \"25\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"address\": \"0x83a7663b2b9d6d3f377a41d03b03ba0021e2f831\",\n" +
	"                        \"currency\": \"USDT\",\n" +
	"                        \"exchange_rate\": \"0.05\",\n" +
	"                      },\n" +
	"                    \"amount_sub_total\": \"25\",\n" +
	"                    \"amount_total\": \"25\",\n" +
	"                    \"note\": \"Test profit withdraw to USDT\",\n" +
	"                  }\n" +
	"              blockchain_txn:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"currency\": \"USDT\",\n" +
	"                    \"src_channel_type\": 10003,\n" +
	"                    \"src_channel_amount\": \"1\",\n" +
	"                    \"dst_channel_type\": 100003,\n" +
	"                    \"dst_channel_amount\": \"1\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"to_address\": \"0xf12Db5BeC81DD6f41366EbcC599508212b751479\",\n" +
	"                      },\n" +
	"                    \"amount_sub_total\": \"1\",\n" +
	"                    \"amount_total\": \"1\",\n" +
	"                    \"note\": \"Test submit USDT txn\",\n" +
	"                  }\n" +
	"              torque_purchase:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"currency\": \"USDT\",\n" +
	"                    \"src_channel_type\": 10003,\n" +
	"                    \"src_channel_amount\": \"2\",\n" +
	"                    \"dst_channel_type\": 100004,\n" +
	"                    \"dst_channel_amount\": \"2\",\n" +
	"                  }\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      order_id:\n" +
	"                        $ref: \"#/components/schemas/ID\"\n" +
	"  /v1/payment/order/execute/:\n" +
	"    post:\n" +
	"      summary: Execute an Order\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                order_id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                auth_code:\n" +
	"                  $ref: \"#/components/schemas/AuthCode\"\n" +
	"              required:\n" +
	"                - order_id\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    $ref: \"#/components/schemas/Order\"\n" +
	"  /v1/helper/blockchain/address/validate/:\n" +
	"    post:\n" +
	"      summary: Validate a Blockchain Address\n" +
	"      tags:\n" +
	"        - Helper\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                network:\n" +
	"                  $ref: \"#/components/schemas/Network\"\n" +
	"                address:\n" +
	"                  $ref: \"#/components/schemas/Address\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - address\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      is_valid:\n" +
	"                        type: boolean\n" +
	"                      address:\n" +
	"                        $ref: \"#/components/schemas/Address\"\n" +
	"  /v1/s2s/portfolio/currency/get/:\n" +
	"    post:\n" +
	"      summary: Get User Currency Info\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"              required:\n" +
	"                - uid\n" +
	"                - currency\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      currency:\n" +
	"                        $ref: \"#/components/schemas/Currency\"\n" +
	"                      price_usd:\n" +
	"                        $ref: \"#/components/schemas/Decimal\"\n" +
	"                      price_usdt:\n" +
	"                        $ref: \"#/components/schemas/Decimal\"\n" +
	"                      account_no:\n" +
	"                        $ref: \"#/components/schemas/Address\"\n" +
	"                      balance:\n" +
	"                        $ref: \"#/components/schemas/Decimal\"\n" +
	"  /v1/s2s/payment/order/init/:\n" +
	"    post:\n" +
	"      summary: Init a new Order\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                src_channel_type:\n" +
	"                  $ref: \"#/components/schemas/ChannelType\"\n" +
	"                src_channel_id:\n" +
	"                  type: integer\n" +
	"                src_channel_ref:\n" +
	"                  $ref: \"#/components/schemas/Reference\"\n" +
	"                src_channel_amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                src_channel_context:\n" +
	"                  type: object\n" +
	"                dst_channel_type:\n" +
	"                  $ref: \"#/components/schemas/ChannelType\"\n" +
	"                dst_channel_id:\n" +
	"                  type: integer\n" +
	"                dst_channel_ref:\n" +
	"                  $ref: \"#/components/schemas/Reference\"\n" +
	"                dst_channel_amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                dst_channel_context:\n" +
	"                  type: object\n" +
	"                amount_sub_total:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                amount_total:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                note:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"              required:\n" +
	"                - uid\n" +
	"                - currency\n" +
	"                - src_channel_type\n" +
	"                - src_channel_amount\n" +
	"                - dst_channel_type\n" +
	"                - dst_channel_amount\n" +
	"                - amount_sub_total\n" +
	"                - amount_total\n" +
	"            examples:\n" +
	"              transfer_p2p:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"uid\": 1,\n" +
	"                    \"currency\": \"TORQ\",\n" +
	"                    \"src_channel_type\": 10000,\n" +
	"                    \"src_channel_amount\": \"12.34567898\",\n" +
	"                    \"dst_channel_type\": 100002,\n" +
	"                    \"dst_channel_amount\": \"12.34567898\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"note\": \"Hello World!\",\n" +
	"                        \"user_identity\": \"powerman\",\n" +
	"                      },\n" +
	"                    \"amount_sub_total\": \"12.34567898\",\n" +
	"                    \"amount_total\": \"12.34567898\",\n" +
	"                    \"note\": \"Test transfer\",\n" +
	"                  }\n" +
	"              profit_reinvest:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"uid\": 1,\n" +
	"                    \"currency\": \"TORQ\",\n" +
	"                    \"src_channel_type\": 10000,\n" +
	"                    \"src_channel_amount\": \"20\",\n" +
	"                    \"dst_channel_type\": 100100,\n" +
	"                    \"dst_channel_amount\": \"20\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"user_identity\": \"powerman\",\n" +
	"                        \"currency\": \"USDT\",\n" +
	"                        \"exchange_rate\": \"0.05\",\n" +
	"                      },\n" +
	"                    \"amount_sub_total\": \"20\",\n" +
	"                    \"amount_total\": \"20\",\n" +
	"                    \"note\": \"Test reinvest USDT\",\n" +
	"                  }\n" +
	"              profit_withdraw:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"uid\": 1,\n" +
	"                    \"currency\": \"TORQ\",\n" +
	"                    \"src_channel_type\": 10000,\n" +
	"                    \"src_channel_amount\": \"25\",\n" +
	"                    \"dst_channel_type\": 100101,\n" +
	"                    \"dst_channel_amount\": \"25\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"address\": \"0x83a7663b2b9d6d3f377a41d03b03ba0021e2f831\",\n" +
	"                        \"currency\": \"USDT\",\n" +
	"                        \"exchange_rate\": \"0.05\",\n" +
	"                      },\n" +
	"                    \"amount_sub_total\": \"25\",\n" +
	"                    \"amount_total\": \"25\",\n" +
	"                    \"note\": \"Test profit withdraw to USDT\",\n" +
	"                  }\n" +
	"              blockchain_txn:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"uid\": 1,\n" +
	"                    \"currency\": \"USDT\",\n" +
	"                    \"src_channel_type\": 10003,\n" +
	"                    \"src_channel_amount\": \"1\",\n" +
	"                    \"dst_channel_type\": 100003,\n" +
	"                    \"dst_channel_amount\": \"1\",\n" +
	"                    \"dst_channel_context\":\n" +
	"                      {\n" +
	"                        \"to_address\": \"0xf12Db5BeC81DD6f41366EbcC599508212b751479\",\n" +
	"                      },\n" +
	"                    \"amount_sub_total\": \"1\",\n" +
	"                    \"amount_total\": \"1\",\n" +
	"                    \"note\": \"Test submit USDT txn\",\n" +
	"                  }\n" +
	"              torque_purchase:\n" +
	"                value:\n" +
	"                  {\n" +
	"                    \"uid\": 1,\n" +
	"                    \"currency\": \"USDT\",\n" +
	"                    \"src_channel_type\": 10003,\n" +
	"                    \"src_channel_amount\": \"2\",\n" +
	"                    \"dst_channel_type\": 100004,\n" +
	"                    \"dst_channel_amount\": \"2\",\n" +
	"                  }\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      order_id:\n" +
	"                        $ref: \"#/components/schemas/ID\"\n" +
	"  /v1/s2s/payment/order/execute/:\n" +
	"    post:\n" +
	"      summary: Execute an Order\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                order_id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"              required:\n" +
	"                - uid\n" +
	"                - order_id\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    $ref: \"#/components/schemas/Order\"\n" +
	"  /v1/s2s/helper/blockchain/address/validate/:\n" +
	"    post:\n" +
	"      summary: Validate a Blockchain Address\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"        - Helper\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                network:\n" +
	"                  $ref: \"#/components/schemas/Network\"\n" +
	"                address:\n" +
	"                  $ref: \"#/components/schemas/Address\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - address\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      is_valid:\n" +
	"                        type: boolean\n" +
	"                      address:\n" +
	"                        $ref: \"#/components/schemas/Address\"\n" +
	"  /v1/s2s/meta/handshake/:\n" +
	"    post:\n" +
	"      summary: Init meta data\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"        - Meta\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      features:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          type: object\n" +
	"                          properties:\n" +
	"                            code:\n" +
	"                              type: string\n" +
	"                              example: kyc\n" +
	"                            name:\n" +
	"                              type: string\n" +
	"                              example: KYC\n" +
	"                            is_available:\n" +
	"                              type: boolean\n" +
	"                      blockchain_networks:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          type: object\n" +
	"                          properties:\n" +
	"                            code:\n" +
	"                              $ref: \"#/components/schemas/Network\"\n" +
	"                            currency:\n" +
	"                              $ref: \"#/components/schemas/Currency\"\n" +
	"                            name:\n" +
	"                              type: string\n" +
	"                              example: Ethereum\n" +
	"                            token_transfer_code_name:\n" +
	"                              type: string\n" +
	"                              example: ERC20\n" +
	"                      network_currencies:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          type: object\n" +
	"                          properties:\n" +
	"                            currency:\n" +
	"                              $ref: \"#/components/schemas/Currency\"\n" +
	"                            network:\n" +
	"                              $ref: \"#/components/schemas/Network\"\n" +
	"                            priority:\n" +
	"                              type: integer\n" +
	"                              example: 2\n" +
	"                            withdrawal_fee:\n" +
	"                              $ref: \"#/components/schemas/Decimal\"\n" +
	"  /v1/s2s/meta/currency-info/:\n" +
	"    post:\n" +
	"      summary: Get system currency info\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"        - Meta\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      currencies:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          $ref: \"#/components/schemas/CurrencyInfo\"\n" +
	"  /v1/s2s/risk/action-lock/get/:\n" +
	"    post:\n" +
	"      summary: Checkout status user action\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                action_type:\n" +
	"                  $ref: \"#/components/schemas/LockAction\"\n" +
	"              required:\n" +
	"                - uid\n" +
	"                - action_type\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      is_locked:\n" +
	"                        type: boolean\n" +
	""
